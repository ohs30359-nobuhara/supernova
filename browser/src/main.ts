import {HeadlessBrowser} from "./headlessBrowser";
import cac, {CAC} from 'cac';
import {createFileWithDirectory} from "./file";

(async () => main())();

async function main(): Promise<void> {
  const cli: CAC = cac();
  const browser: HeadlessBrowser = await HeadlessBrowser.New();
  try {
    const parsed = cli
      .option("--performance <performance>", "Get core web vital. Please specify file name.")
      .option("--screenshot <screenshot>", "Take a screenshot. Please specify file name.")
      .parse();

    if (parsed.args.length === 0) {
      throw Error();
    }
    const url: string = parsed.args[0];

    const options: {[k: string]: any} = parsed.options
    const performance: string | undefined = options["performance"];
    const screenshot: string | undefined = options["screenshot"];

    await browser.move(url)

    if (screenshot) {
      await browser.screenshot(screenshot);
    }

    if (performance) {
      const {html, json} = await browser.coreWebVital();
      createFileWithDirectory(`${performance}.html`, html);
      createFileWithDirectory(`${performance}.json`,JSON.stringify(json));
    }

    process.exit(0)
  } catch (e) {
    console.error(e);

    process.exit(1)
  } finally {
    browser.kill().catch(e => console.error(e));
  }
}
