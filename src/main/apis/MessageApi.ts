import { logger } from '../../utils/logger';

export class MessageApi {
  asyncMessage() {
    return "async message from main";
  }

  syncMessage(arg: any) {
    logger.info(arg);
    return 'main-process pong message';
  }
}
