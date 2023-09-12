
import { StaticServerUrl } from '../../wailsjs/go/apis/StaticInfos';

export const StaticDictServerURL = function (): Promise<string> {
    return StaticServerUrl()
}
