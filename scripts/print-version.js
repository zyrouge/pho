const p = require("path");
const fs = require("fs/promises");

const rootDir = p.resolve(__dirname, "..");
const metaFile = p.join(rootDir, "core/meta.go");

const start = async () => {
    const content = await fs.readFile(metaFile, "utf-8");
    const version = content.match(/AppVersion = "(\d+.\d+.\d+)"/)[1];
    console.log(version);
};

start();
