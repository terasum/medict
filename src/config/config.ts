import { app } from 'electron';
import path from 'path';
import fs from 'fs';

export function getResourceRootPath() {
  const userResourcePath = app.getPath('userData');
  const resourceRootPath = path.resolve(userResourcePath, 'resources', 'cache');
  if (!fs.existsSync(resourceRootPath)) {
    console.info(`create new directory ${resourceRootPath}`);
    fs.mkdirSync(resourceRootPath, { recursive: true });
  }
  return resourceRootPath;
}

export function getConfigJsonPath() {
  const userResourcePath = app.getPath('userData');
  const userConfigDir = path.resolve(userResourcePath, 'config');
  if (!fs.existsSync(userConfigDir)) {
    console.info(`create new directory ${userConfigDir}`);
    fs.mkdirSync(userConfigDir, { recursive: true });
  }

  const configFilePath = path.resolve(
    userResourcePath,
    'config',
    'medict.json'
  );
  console.info('configFilePath %s', configFilePath);
  if (!fs.existsSync(configFilePath)) {
    fs.writeFileSync(configFilePath, '{"dicts":[]}');
    console.info('write configFilePath %s', configFilePath);
  }
  return configFilePath;
}

export function getLoggerFilePath() {
  const userResourcePath = app.getPath('userData');
  const userLoggerDir = path.resolve(userResourcePath, 'logger');
  if (!fs.existsSync(userLoggerDir)) {
    console.info(`create new directory ${userLoggerDir}`);
    fs.mkdirSync(userLoggerDir, { recursive: true });
  }

  const loggerFilePath = path.resolve(userResourcePath, 'logger', 'medict.log');
  console.info('logfile path %s', loggerFilePath);
  return loggerFilePath;
}
