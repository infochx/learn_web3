const { extendEnvironment } = require("hardhat/config");

require("@nomicfoundation/hardhat-toolbox");
require("dotenv").config();
require("@openzeppelin/hardhat-upgrades");

// 任务：打印所有账户地址
task("accounts", "Prints the list of accounts", async (taskArgs, hre) => {
  const accounts = await hre.ethers.getSigners();

  for (const account of accounts) {
    console.log(account.address);
  }
});

// 扩展 Hardhat 环境，添加 Web3 实例
// extendEnvironment((hre) => {
//   const Web3 = require("web3");
//   hre.Web3 = Web3;
//   hre.web3 = new Web3(hre.ethers.provider);
// });

module.exports = {
  solidity: {
    version: "0.8.28", // 确保与合约 pragma 版本一致
    settings: {
      optimizer: { enabled: true, runs: 200 }, // 启用优化，加快部署
    },
  },

  networks: {
    hardhat: {
      blockGasLimit: 30000000, // 足够大的 Gas 限制，避免部署因 Gas 不足失败
      gasPrice: 8000000000, // 合理 Gas 价格，避免 pending
    },
    sepolia: {
      url: process.env.SEPOLIA_RPC_URL || "",
      accounts: process.env.PRIVATE_KEY ? [process.env.PRIVATE_KEY] : [],
    },
  },

  etherscan: {
    apiKey: process.env.ETHERSCAN_API_KEY,
  },
};
