import { app } from 'electron';
import path from 'path';
import fs from 'fs';
import { logger } from '../utils/logger';

export function getResourceRootPath() {
  const userResourcePath = app.getPath('userData');
  const resourceRootPath = path.resolve(userResourcePath, 'resources', 'cache');
  if (!fs.existsSync(resourceRootPath)) {
    logger.info(`create new directory ${resourceRootPath}`);
    fs.mkdirSync(resourceRootPath, { recursive: true });
  }
  return resourceRootPath;
}

export function getConfigJsonPath() {
  const userResourcePath = app.getPath('userData');
  const userConfigDir = path.resolve(userResourcePath, 'config');
  if (!fs.existsSync(userConfigDir)) {
    logger.info(`create new directory ${userConfigDir}`);
    fs.mkdirSync(userConfigDir, { recursive: true });
  }

  const configFilePath = path.resolve(
    userResourcePath,
    'config',
    'medict.json'
  );
  logger.info('configFilePath %s', configFilePath);
  if (!fs.existsSync(configFilePath)) {
    fs.writeFileSync(configFilePath, '{"dicts":[]}');
    logger.info('write configFilePath %s', configFilePath);
  }
  return configFilePath;
}
