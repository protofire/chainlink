let LINK = artifacts.require("../../../solidity/contracts/lib/LinkToken.sol");

module.exports = function(deployer) {
  deployer.deploy(LINK);
};

