
import path from 'path';
import fs from 'fs';
import { preloadContent } from '../../renderer/preload/webview.preload';
import { logger } from '../../utils/logger';
import { SyncMainAPI } from '../../main/rpc.sync.main.reference';

let __globalConfig__: Configuration;

export default class Configuration {
    userResourcePath: string;
    resourceRootPath: string;
    userConfigDir: string;
    configJsonFilePath: string;
    userLoggerDir: string;
    loggerFilePath: string;
    preloadFilePath: string;
    preloadDir: string;

    public static newInstance() {
        if (!__globalConfig__) {
            __globalConfig__ = new Configuration();
        }
        return __globalConfig__;
    }

    private constructor() {
        this.userResourcePath = SyncMainAPI.syncGetResourceRootPath();
        this.resourceRootPath = path.resolve(this.userResourcePath, 'resources', 'cache');
        this.userConfigDir = path.resolve(this.userResourcePath, 'config');
        this.configJsonFilePath = path.resolve(this.userResourcePath, 'config', 'medict.json');
        this.userLoggerDir = path.resolve(this.userResourcePath, 'logger');
        this.loggerFilePath = path.resolve(this.userResourcePath, 'logger', 'medict.log');
        this.preloadDir = path.resolve(this.userResourcePath, 'preload');
        this.preloadFilePath = path.resolve(this.userResourcePath, 'preload', 'medict-webview.preload.js');

        if (!fs.existsSync(this.resourceRootPath)) {
            console.info(`create new directory ${this.resourceRootPath}`);
            fs.mkdirSync(this.resourceRootPath, { recursive: true });
        }

        if (!fs.existsSync(this.userConfigDir)) {
            console.info(`create new directory ${this.userConfigDir}`);
            fs.mkdirSync(this.userConfigDir, { recursive: true });
        }

        if (!fs.existsSync(this.configJsonFilePath)) {
            console.info('write default config FilePath %s', this.configJsonFilePath);
            fs.writeFileSync(this.configJsonFilePath, '{"dicts":[]}');
            console.info('write configFilePath %s', this.configJsonFilePath);
        }

        if (!fs.existsSync(this.userLoggerDir)) {
            console.info(`create new directory ${this.userLoggerDir}`);
            fs.mkdirSync(this.userLoggerDir, { recursive: true });
        }


        if (!fs.existsSync(this.preloadDir)) {
            console.info(`create new directory ${this.preloadDir}`);
            fs.mkdirSync(this.preloadDir, { recursive: true });
        }

        // write preload file
        fs.writeFile(this.preloadFilePath, preloadContent, () => {
            logger.info(
                `write ${this.preloadFilePath} successfully`
            );
        })

    }


    getResourceRootPath() {
        return this.resourceRootPath;
    }

    getConfigJsonPath() {
        return this.configJsonFilePath;
    }

    getLoggerFilePath() {
        return this.loggerFilePath;
    }

    webviewPreloadFilePath() {
        return this.preloadFilePath;
    }
    getYoudaoAppID() {
        return '';
    }
    getYoudaoAppKey() {
        return ''
    }
    getBaiduAppID() {
        return '';
    }
    getBaiduAppKey() {
        return ''
    }
}
