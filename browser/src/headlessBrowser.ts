import puppeteer, {Browser, Page} from "puppeteer";
import {convertJson} from "./performance";
import {Performance} from "./performance"
import {createFileWithDirectory} from "./file";

export class HeadlessBrowser {
  private readonly browser: Browser
  private activePage: Page
  /**
   * @constructor
   * @param browser
   * @param page
   */
  constructor(browser: Browser, page: Page) {
    this.browser = browser;
    this.activePage = page;
  }

  public static async New(): Promise<HeadlessBrowser> {
    const b: Browser = await puppeteer.launch({
      headless: false,
      defaultViewport: null,
      ignoreDefaultArgs: ['--enable-automation']
    });
    const p: Page = await b.newPage();

    return new HeadlessBrowser(b, p);
  }

  /**
   * ページ遷移
   * @param url
   */
  public async move(url: string) {
    this.activePage = await this.browser.newPage();
    await this.activePage.goto(url, {
      waitUntil: "domcontentloaded",
    });
  }

  /**
   * スクリーンショットの取得
   */
  public async screenshot(path: string): Promise<void> {
    const buf: Buffer = await this.activePage.screenshot();
    createFileWithDirectory(path, buf.toString());
  }

  /**
   * CoreWebVitalの取得
   */
  public async coreWebVital(): Promise<{html: string, json: Performance[]}> {
    // puppeteerからでは開けないので chromeから開く
    // https://github.com/GoogleChrome/lighthouse/issues/15124
    const lighthouse = require('lighthouse/core/index.cjs');
    const result = await lighthouse(this.activePage.url(), {
      logLevel: "error",
      output: "html",
      port: + new URL(this.browser.wsEndpoint()).port
    });

    if (!result) {
      throw new Error("failed to retrieve report");
    }

    return {
      html: result.report as string,
      json: convertJson(result)
    }
  }

  /**
   * ブラウザを閉じる
   */
  public async kill(): Promise<void> {
    await this.browser.close()
  }
}
