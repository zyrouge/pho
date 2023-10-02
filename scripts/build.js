const p = require("path");
const { spawnSync } = require("child_process");

const rootDir = p.resolve(__dirname, "..");
const outputDir = p.join(rootDir, "dist");

const buildArchs = ["386", "amd64", "arm", "arm64"];

const start = async () => {
    for (const arch of buildArchs) {
        const cwd = process.cwd;
        const env = {
            ...process.env,
            GOOS: "linux",
            GOARCH: arch,
        };
        const outputFile = p.join(outputDir, `pho-${arch}`);
        console.log(`Building ${outputFile}...`);
        const result = spawnSync(
            "go",
            ["build", '-ldflags="-s -w"', "-o", outputFile],
            {
                cwd,
                env,
            }
        );
        if (result.signal !== 0) {
            throw new Error("Build failed");
        }
    }
};

start();
