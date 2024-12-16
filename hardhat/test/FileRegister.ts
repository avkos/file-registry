import { expect } from "chai";
import hre from "hardhat";
import {v4} from "uuid";
import {FileRegistry} from "../typechain-types";

describe("FileRegister", function () {
  async function deploy() {
    const [owner] = await hre.ethers.getSigners();
    const FileRegistry = await hre.ethers.getContractFactory("FileRegistry");
    const fileRegistry = await FileRegistry.deploy();
    return { fileRegistry, owner };
  }

  describe("Deployment", function () {
    let contract:FileRegistry
    before(async()=>{
      const {fileRegistry} = await deploy()
      contract=fileRegistry
    })

    it("Should deploy contract", async function () {
      const address = await contract.getAddress()
      expect(address).to.properAddress;
      expect(address.length).to.equal(42);
    });

    it("Should save a filePath", async function () {
      const res = await contract.save('some/file/path', v4())
      expect(res.hash).to.exist
      expect(res.hash.length).to.equal(66);
    });

    it("Should return cid for a filePath", async function () {
      const cid = v4()
      const path = `some/file/path/${v4()}`
      await contract.save(path, cid)
      const resCid = await contract.get(path)
      expect(resCid).to.equal(cid);
    });

    it("Should replace a cid", async function () {
      const cid = v4()
      const replacedCid = v4()
      const path = `some/file/path/check/replace`
      await contract.save(path, cid)
      const resCid = await contract.get(path)
      expect(resCid).to.equal(cid);
      await contract.save(path, replacedCid)
      const resReplacedCid = await contract.get(path)
      expect(resReplacedCid).to.equal(replacedCid);

    });

    it("Should return an empty string if CID was not set for a file path", async function () {
      const unknownFilePath = "/unknown/path.txt";
      const storedCid = await contract.get(unknownFilePath);
      expect(storedCid).to.equal("");
    });

  });
});
