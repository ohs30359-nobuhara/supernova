import puppeteer, {Browser, Page} from "puppeteer";
import {writeFileSync} from "fs";
import {convertJson} from "./performance";
import {Performance} from "./performance"

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
    writeFileSync(path, buf);
  }

  /**
   * CoreWebVitalの取得
   */
  public async coreWebVital(): Promise<Performance[]> {
    // puppeteerからでは開けないので chromeから開く
    // https://github.com/GoogleChrome/lighthouse/issues/15124
    const lighthouse = require('lighthouse/core/index.cjs');
    const report = await lighthouse(this.activePage.url(), {
      logLevel: "error",
      output: "json",
      port: + new URL(this.browser.wsEndpoint()).port
    });

    if (!report) {
      throw new Error("failed to retrieve report");
    }

    return convertJson(report);
  }

  /**
   * ブラウザを閉じる
   */
  public async kill(): Promise<void> {
    await this.browser.close()
  }
}
