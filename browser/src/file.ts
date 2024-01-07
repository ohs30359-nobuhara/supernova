import {dirname} from "path";
import {existsSync, mkdirSync, writeFileSync} from "fs";

/**
 * Dirを考慮してファイルを作成する
 * @param filePath
 * @param fileContent
 */
export function createFileWithDirectory(filePath: string, fileContent: string): void {
  const directoryPath = dirname(filePath);

  // ディレクトリが存在するか確認
  if (!existsSync(directoryPath)) {
    // ディレクトリを作成
    mkdirSync(directoryPath, { recursive: true });
  }

  // ファイルを作成
  writeFileSync(filePath, fileContent);
}
