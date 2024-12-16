const { ethers } = require("hardhat");

async function main() {
    const FileRegistry = await ethers.getContractFactory("FileRegistry");
    const registry = await FileRegistry.deploy();
    // get private key from hardhat config

    const contractAddress = await registry.getAddress()
    console.log("FileRegistry deployed to:", contractAddress)
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});