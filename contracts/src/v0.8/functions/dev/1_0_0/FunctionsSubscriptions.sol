// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IFunctionsSubscriptions} from "./interfaces/IFunctionsSubscriptions.sol";
import {ERC677ReceiverInterface} from "../../../interfaces/ERC677ReceiverInterface.sol";
import {LinkTokenInterface} from "../../../interfaces/LinkTokenInterface.sol";
import {IFunctionsBilling} from "./interfaces/IFunctionsBilling.sol";
import {IFunctionsRequest} from "./interfaces/IFunctionsRequest.sol";
import {IFunctionsRouter} from "./interfaces/IFunctionsRouter.sol";
import {SafeCast} from "../../../vendor/openzeppelin-solidity/v4.8.0/contracts/utils/SafeCast.sol";

/**
 * @title Functions Subscriptions contract
 * @notice Contract that coordinates payment from users to the nodes of the Decentralized Oracle Network (DON).
 * @dev THIS CONTRACT HAS NOT GONE THROUGH ANY SECURITY REVIEW. DO NOT USE IN PROD.
 */
abstract contract FunctionsSubscriptions is IFunctionsSubscriptions, ERC677ReceiverInterface {
  // ================================================================
  // |                      Subscription state                      |
  // ================================================================

  // We make the sub count public so that its possible to
  // get all the current subscriptions via getSubscription.
  uint64 private s_currentSubscriptionId;

  // s_totalBalance tracks the total LINK sent to/from
  // this contract through onTokenTransfer, cancelSubscription and oracleWithdraw.
  // A discrepancy with this contract's LINK balance indicates that someone
  // sent tokens using transfer and so we may need to use recoverFunds.
  uint96 private s_totalBalance;

  // link token address
  LinkTokenInterface private s_linkToken;

  mapping(uint64 subscriptionId => IFunctionsSubscriptions.Subscription) internal s_subscriptions;
  mapping(address consumer => mapping(uint64 subscriptionId => IFunctionsSubscriptions.Consumer)) internal s_consumers;

  event SubscriptionCreated(uint64 indexed subscriptionId, address owner);
  event SubscriptionFunded(uint64 indexed subscriptionId, uint256 oldBalance, uint256 newBalance);
  event SubscriptionConsumerAdded(uint64 indexed subscriptionId, address consumer);
  event SubscriptionConsumerRemoved(uint64 indexed subscriptionId, address consumer);
  event SubscriptionCanceled(uint64 indexed subscriptionId, address fundsRecipient, uint256 fundsAmount);
  event SubscriptionOwnerTransferRequested(uint64 indexed subscriptionId, address from, address to);
  event SubscriptionOwnerTransferred(uint64 indexed subscriptionId, address from, address to);

  error TooManyConsumers();
  error InsufficientBalance();
  error InvalidConsumer();
  error ConsumerRequestsInFlight();
  error InvalidSubscription();
  error OnlyCallableFromLink();
  error InvalidCalldata();
  error MustBeSubscriptionOwner();
  error PendingRequestExists();
  error MustBeProposedOwner();
  error BalanceInvariantViolated(uint256 internalBalance, uint256 externalBalance); // Should never happen
  event FundsRecovered(address to, uint256 amount);

  // @dev NOP balances are held as a single amount. The breakdown is held by the Coordinator.
  mapping(address coordinator => uint96 balanceJuelsLink) private s_withdrawableTokens;

  // ================================================================
  // |                       Request state                          |
  // ================================================================

  mapping(bytes32 requestId => bytes32 commitmentHash) internal s_requestCommitments;

  struct Receipt {
    uint96 callbackGasCostJuels;
    uint96 totalCostJuels;
  }

  event RequestTimedOut(bytes32 indexed requestId);

  // ================================================================
  // |                       Initialization                         |
  // ================================================================
  constructor(address link) {
    s_linkToken = LinkTokenInterface(link);
  }

  // ================================================================
  // |                      Getter methods                          |
  // ================================================================
  /**
   * @inheritdoc IFunctionsSubscriptions
   */
  function getTotalBalance() external view override returns (uint96) {
    return s_totalBalance;
  }

  /**
   * @inheritdoc IFunctionsSubscriptions
   */
  function getSubscriptionCount() external view override returns (uint64) {
    return s_currentSubscriptionId;
  }

  /**
   * @inheritdoc IFunctionsSubscriptions
   */
  function getSubscription(uint64 subscriptionId) external view override returns (Subscription memory) {
    _isValidSubscription(subscriptionId);

    return s_subscriptions[subscriptionId];
  }

  /**
   * @inheritdoc IFunctionsSubscriptions
   */
  function getConsumer(
    address client,
    uint64 subscriptionId
  ) external view override returns (bool allowed, uint64 initiatedRequests, uint64 completedRequests) {
    allowed = s_consumers[client][subscriptionId].allowed;
    initiatedRequests = s_consumers[client][subscriptionId].initiatedRequests;
    completedRequests = s_consumers[client][subscriptionId].completedRequests;
  }

  // ================================================================
  // |                      Internal checks                         |
  // ================================================================

  function _isValidSubscription(uint64 subscriptionId) internal view {
    if (s_subscriptions[subscriptionId].owner == address(0)) {
      revert InvalidSubscription();
    }
  }

  function _isValidConsumer(address client, uint64 subscriptionId) internal view {
    if (!s_consumers[client][subscriptionId].allowed) {
      revert InvalidConsumer();
    }
  }

  // ================================================================
  // |                 Internal Payment methods                     |
  // ================================================================
  /**
   * @notice Sets a request as in-flight
   * @dev Only callable within the Router
   */
  function _markRequestInFlight(address client, uint64 subscriptionId, uint96 estimatedTotalCostJuels) internal {
    // Earmark subscription funds
    s_subscriptions[subscriptionId].blockedBalance += estimatedTotalCostJuels;

    // Increment sent requests
    s_consumers[client][subscriptionId].initiatedRequests += 1;
  }

  /**
   * @notice Moves funds from one subscription account to another.
   * @dev Only callable by the Coordinator contract that is saved in the request commitment
   */
  function _pay(
    uint64 subscriptionId,
    uint96 estimatedTotalCostJuels,
    address client,
    uint96 adminFee,
    uint96 juelsPerGas,
    uint96 gasUsed,
    uint96 costWithoutCallbackJuels
  ) internal returns (Receipt memory receipt) {
    uint96 callbackGasCostJuels = juelsPerGas * gasUsed;
    uint96 totalCostJuels = costWithoutCallbackJuels + adminFee + callbackGasCostJuels;

    receipt = Receipt(callbackGasCostJuels, totalCostJuels);

    // Charge the subscription
    s_subscriptions[subscriptionId].balance -= totalCostJuels;

    // Pay the DON's fees and gas reimbursement
    s_withdrawableTokens[msg.sender] += costWithoutCallbackJuels + callbackGasCostJuels;

    // Pay out the administration fee
    s_withdrawableTokens[address(this)] += adminFee;

    // Unblock earmarked funds
    s_subscriptions[subscriptionId].blockedBalance -= estimatedTotalCostJuels;
    // Increment finished requests
    s_consumers[client][subscriptionId].completedRequests += 1;
  }

  // ================================================================
  // |                      Owner methods                           |
  // ================================================================
  /**
   * @inheritdoc IFunctionsSubscriptions
   */
  function ownerCancelSubscription(uint64 subscriptionId) external override {
    _onlyRouterOwner();
    _isValidSubscription(subscriptionId);
    _cancelSubscriptionHelper(subscriptionId, s_subscriptions[subscriptionId].owner);
  }

  /**
   * @inheritdoc IFunctionsSubscriptions
   */
  function recoverFunds(address to) external override {
    _onlyRouterOwner();
    uint256 externalBalance = s_linkToken.balanceOf(address(this));
    uint256 internalBalance = uint256(s_totalBalance);
    if (internalBalance > externalBalance) {
      revert BalanceInvariantViolated(internalBalance, externalBalance);
    }
    if (internalBalance < externalBalance) {
      uint256 amount = externalBalance - internalBalance;
      s_linkToken.transfer(to, amount);
      emit FundsRecovered(to, amount);
    }
    // If the balances are equal, nothing to be done.
  }

  /**
   * @notice Owner withdraw LINK earned through admin fees
   * @notice If amount is 0 the full balance will be withdrawn
   * @param recipient where to send the funds
   * @param amount amount to withdraw
   */
  function ownerWithdraw(address recipient, uint96 amount) external {
    _onlyRouterOwner();
    if (amount == 0) {
      amount = s_withdrawableTokens[address(this)];
    }
    if (s_withdrawableTokens[address(this)] < amount) {
      revert InsufficientBalance();
    }
    s_withdrawableTokens[address(this)] -= amount;
    s_totalBalance -= amount;
    if (!s_linkToken.transfer(recipient, amount)) {
      uint256 externalBalance = s_linkToken.balanceOf(address(this));
      uint256 internalBalance = uint256(s_totalBalance);
      revert BalanceInvariantViolated(internalBalance, externalBalance);
    }
  }

  // ================================================================
  // |                     Coordinator methods                      |
  // ================================================================
  /**
   * @inheritdoc IFunctionsSubscriptions
   */
  function oracleWithdraw(address recipient, uint96 amount) external override {
    _whenNotPaused();
    if (amount == 0) {
      revert InvalidCalldata();
    }
    if (s_withdrawableTokens[msg.sender] < amount) {
      revert InsufficientBalance();
    }
    s_withdrawableTokens[msg.sender] -= amount;
    s_totalBalance -= amount;
    if (!s_linkToken.transfer(recipient, amount)) {
      revert InsufficientBalance();
    }
  }

  // ================================================================
  // |                   Deposit helper method                      |
  // ================================================================
  function onTokenTransfer(address /* sender */, uint256 amount, bytes calldata data) external override {
    _whenNotPaused();
    if (msg.sender != address(s_linkToken)) {
      revert OnlyCallableFromLink();
    }
    if (data.length != 32) {
      revert InvalidCalldata();
    }
    uint64 subscriptionId = abi.decode(data, (uint64));
    if (s_subscriptions[subscriptionId].owner == address(0)) {
      revert InvalidSubscription();
    }
    // We do not check that the msg.sender is the subscription owner,
    // anyone can fund a subscription.
    uint256 oldBalance = s_subscriptions[subscriptionId].balance;
    s_subscriptions[subscriptionId].balance += uint96(amount);
    s_totalBalance += uint96(amount);
    emit SubscriptionFunded(subscriptionId, oldBalance, oldBalance + amount);
  }

  // ================================================================
  // |                    Subscription methods                      |
  // ================================================================
  /**
   * @inheritdoc IFunctionsSubscriptions
   */
  function createSubscription() external override returns (uint64 subscriptionId) {
    _whenNotPaused();
    _onlySenderThatAcceptedToS();
    subscriptionId = ++s_currentSubscriptionId;
    s_subscriptions[subscriptionId] = Subscription({
      balance: 0,
      blockedBalance: 0,
      owner: msg.sender,
      requestedOwner: address(0),
      consumers: new address[](0),
      flags: bytes32(0)
    });

    emit SubscriptionCreated(subscriptionId, msg.sender);
  }

  /**
   * @inheritdoc IFunctionsSubscriptions
   */
  function proposeSubscriptionOwnerTransfer(uint64 subscriptionId, address newOwner) external override {
    _whenNotPaused();
    _onlySubscriptionOwner(subscriptionId);
    _onlySenderThatAcceptedToS();

    // Proposing to address(0) would never be claimable, so don't need to check.

    if (s_subscriptions[subscriptionId].requestedOwner != newOwner) {
      s_subscriptions[subscriptionId].requestedOwner = newOwner;
      emit SubscriptionOwnerTransferRequested(subscriptionId, msg.sender, newOwner);
    }
  }

  /**
   * @inheritdoc IFunctionsSubscriptions
   */
  function acceptSubscriptionOwnerTransfer(uint64 subscriptionId) external override {
    _whenNotPaused();
    _onlySenderThatAcceptedToS();
    address previousOwner = s_subscriptions[subscriptionId].owner;
    address nextOwner = s_subscriptions[subscriptionId].requestedOwner;
    if (nextOwner != msg.sender) {
      revert MustBeProposedOwner();
    }
    s_subscriptions[subscriptionId].owner = msg.sender;
    s_subscriptions[subscriptionId].requestedOwner = address(0);
    emit SubscriptionOwnerTransferred(subscriptionId, previousOwner, msg.sender);
  }

  /**
   * @inheritdoc IFunctionsSubscriptions
   */
  function removeConsumer(uint64 subscriptionId, address consumer) external override {
    _whenNotPaused();
    _onlySubscriptionOwner(subscriptionId);
    _onlySenderThatAcceptedToS();
    Consumer memory consumerData = s_consumers[consumer][subscriptionId];
    if (!consumerData.allowed) {
      revert InvalidConsumer();
    }
    if (consumerData.initiatedRequests != consumerData.completedRequests) {
      revert ConsumerRequestsInFlight();
    }
    // Note bounded by config.maxConsumers
    address[] memory consumers = s_subscriptions[subscriptionId].consumers;
    uint256 lastConsumerIndex = consumers.length - 1;
    for (uint256 i = 0; i < consumers.length; ++i) {
      if (consumers[i] == consumer) {
        address last = consumers[lastConsumerIndex];
        // Storage write to preserve last element
        s_subscriptions[subscriptionId].consumers[i] = last;
        // Storage remove last element
        s_subscriptions[subscriptionId].consumers.pop();
        break;
      }
    }
    delete s_consumers[consumer][subscriptionId];
    emit SubscriptionConsumerRemoved(subscriptionId, consumer);
  }

  /**
   * @inheritdoc IFunctionsSubscriptions
   */
  function addConsumer(uint64 subscriptionId, address consumer) external override {
    _whenNotPaused();
    _onlySubscriptionOwner(subscriptionId);
    _onlySenderThatAcceptedToS();
    // Already maxed, cannot add any more consumers.
    if (s_subscriptions[subscriptionId].consumers.length == _getMaxConsumers()) {
      revert TooManyConsumers();
    }
    if (s_consumers[consumer][subscriptionId].allowed) {
      // Idempotence - do nothing if already added.
      // Ensures uniqueness in s_subscriptions[subscriptionId].consumers.
      return;
    }
    s_consumers[consumer][subscriptionId].allowed = true;
    s_subscriptions[subscriptionId].consumers.push(consumer);

    emit SubscriptionConsumerAdded(subscriptionId, consumer);
  }

  /**
   * @inheritdoc IFunctionsSubscriptions
   */
  function cancelSubscription(uint64 subscriptionId, address to) external override {
    _whenNotPaused();
    _onlySubscriptionOwner(subscriptionId);
    _onlySenderThatAcceptedToS();
    if (_pendingRequestExists(subscriptionId)) {
      revert PendingRequestExists();
    }
    _cancelSubscriptionHelper(subscriptionId, to);
  }

  function _cancelSubscriptionHelper(uint64 subscriptionId, address to) private {
    Subscription memory sub = s_subscriptions[subscriptionId];
    uint96 balance = sub.balance;
    // Note bounded by config.maxConsumers
    // If no consumers, does nothing.
    for (uint256 i = 0; i < sub.consumers.length; ++i) {
      delete s_consumers[sub.consumers[i]][subscriptionId];
    }
    delete s_subscriptions[subscriptionId];
    s_totalBalance -= balance;
    if (!s_linkToken.transfer(to, uint256(balance))) {
      revert InsufficientBalance();
    }
    emit SubscriptionCanceled(subscriptionId, to, balance);
  }

  /**
   * @inheritdoc IFunctionsSubscriptions
   */
  function pendingRequestExists(uint64 subscriptionId) external view override returns (bool) {
    return _pendingRequestExists(subscriptionId);
  }

  function _pendingRequestExists(uint64 subscriptionId) internal view returns (bool) {
    address[] memory consumers = s_subscriptions[subscriptionId].consumers;
    // Iterations will not exceed config.maxConsumers
    for (uint256 i = 0; i < consumers.length; ++i) {
      Consumer memory consumer = s_consumers[consumers[i]][subscriptionId];
      if (consumer.initiatedRequests != consumer.completedRequests) {
        return true;
      }
    }
    return false;
  }

  /**
   * @inheritdoc IFunctionsSubscriptions
   */
  function setFlags(uint64 subscriptionId, bytes32 flags) external override {
    _onlyRouterOwner();
    _isValidSubscription(subscriptionId);
    s_subscriptions[subscriptionId].flags = flags;
  }

  /**
   * @inheritdoc IFunctionsSubscriptions
   */
  function getFlags(uint64 subscriptionId) public view returns (bytes32) {
    return s_subscriptions[subscriptionId].flags;
  }

  function _getMaxConsumers() internal view virtual returns (uint16);

  // ================================================================
  // |                  Request Timeout Methods                     |
  // ================================================================
  /**
   * @inheritdoc IFunctionsSubscriptions
   */
  function timeoutRequests(IFunctionsRequest.Commitment[] calldata requestsToTimeoutByCommitment) external override {
    _whenNotPaused();
    for (uint256 i = 0; i < requestsToTimeoutByCommitment.length; ++i) {
      IFunctionsRequest.Commitment memory request = requestsToTimeoutByCommitment[i];
      bytes32 requestId = request.requestId;

      // Check that request ID is valid
      if (keccak256(abi.encode(request)) != s_requestCommitments[requestId]) {
        revert InvalidCalldata();
      }

      // Check that request has exceeded allowed request time
      if (block.timestamp < request.timeoutTimestamp) {
        revert ConsumerRequestsInFlight();
      }

      IFunctionsBilling coordinator = IFunctionsBilling(request.coordinator);
      coordinator.deleteCommitment(requestId);
      // Release blocked balance
      s_subscriptions[request.subscriptionId].blockedBalance -= request.estimatedTotalCostJuels;
      s_consumers[request.client][request.subscriptionId].completedRequests += 1;
      // Delete commitment
      delete s_requestCommitments[requestId];

      emit RequestTimedOut(requestId);
    }
  }

  // ================================================================
  // |                         Modifiers                            |
  // ================================================================

  function _onlySubscriptionOwner(uint64 subscriptionId) internal view {
    address owner = s_subscriptions[subscriptionId].owner;
    if (owner == address(0)) {
      revert InvalidSubscription();
    }
    if (msg.sender != owner) {
      revert MustBeSubscriptionOwner();
    }
  }

  function _onlySenderThatAcceptedToS() internal virtual;

  function _onlyRouterOwner() internal virtual;

  function _whenNotPaused() internal virtual;
}