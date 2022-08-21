import path from 'path';
import fs from 'fs';
import { staticServerPort } from './ResourceStaticService';

export function getUserResourcePath(): string {
  // const userResourcePath = app.getPath('userData');
  // // 词典静态文件存储目录 http 服务起点
  // const resourceRootPath = path.resolve(userResourcePath, 'resources', 'cache');
  // if (!fs.existsSync(resourceRootPath)) {
  //   console.info(`create new directory ${resourceRootPath}`);
  //   fs.mkdirSync(resourceRootPath, { recursive: true });
  // }

  // return userResourcePath
  return "";
}


export function getResourceRootPath(): string {
  // const userResourcePath = app.getPath('userData');
  // const resourceRootPath = path.resolve(userResourcePath, 'resources', 'cache');
  // return resourceRootPath
  return "";
}

export function getUserConfigDir(): string {
  // const userResourcePath = app.getPath('userData');
  // const userConfigDir = path.resolve(userResourcePath, 'config');
  // return userConfigDir
  return "";
}

export function getConfigJsonFilePath(): string {
  // const userResourcePath = app.getPath('userData');
  // const userConfigDir = path.resolve(userResourcePath, 'config');
  // // 用户 json 配置文件
  // const configJsonFilePath = path.resolve(userResourcePath, 'config', 'medict.json');
  // // 用户配置文件夹
  // if (!fs.existsSync(userConfigDir)) {
  //   console.info(`create new directory ${userConfigDir}`);
  //   fs.mkdirSync(userConfigDir, { recursive: true });
  // }

  // if (!fs.existsSync(configJsonFilePath)) {
  //   fs.writeFileSync(configJsonFilePath, '{"dicts":[]}');
  //   console.info('write configFilePath %s', configJsonFilePath);
  // }

  // return configJsonFilePath

  return "";
}

export function getUserLoggerDir(): string {

  // const userResourcePath = app.getPath('userData');
  // const userLoggerDir = path.resolve(userResourcePath, 'logger');
  // return userLoggerDir
  return "";
}

export function getLoggerFilePath(): string {
  // 用户日志文件夹
  const userLoggerDir = getUserLoggerDir();
  // 日志文件路径
  const loggerFilePath = path.resolve(userLoggerDir, 'medict.log');
  if (!fs.existsSync(userLoggerDir)) {
    console.info(`create new directory ${userLoggerDir}`);
    fs.mkdirSync(userLoggerDir, { recursive: true });
  }
  return loggerFilePath
}

export function getStaticServerPort(): number {
  return staticServerPort();
}

