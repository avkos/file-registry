// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

contract FileRegistry {
    mapping(string => string) private fileToCid;

    event FileSaved(string filePath, string cid);

    function save(string memory filePath, string memory cid) public {
        fileToCid[filePath] = cid;
        emit FileSaved(filePath, cid);
    }

    function get(string memory filePath) public view returns (string memory) {
        return fileToCid[filePath];
    }
}