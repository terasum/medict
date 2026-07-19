import { afterEach, describe, expect, it } from 'vitest';

import { getECDICTStatus, installFullECDICT } from './ecdict';

describe('ECDICT desktop API', () => {
  afterEach(() => {
    delete (window as any).go;
  });

  it('loads compact status and installs the full edition', async () => {
    (window as any).go = {
      main: {
        App: {
          ECDICTStatus: async () => ({ code: 200, data: { entryCount: 50000, edition: 'compact' } }),
          InstallFullECDICT: async () => ({ code: 200, data: { entryCount: 760000, edition: 'full' } }),
        },
      },
    };

    await expect(getECDICTStatus()).resolves.toEqual({ entryCount: 50000, edition: 'compact' });
    await expect(installFullECDICT()).resolves.toEqual({ entryCount: 760000, edition: 'full' });
  });

  it('surfaces backend errors', async () => {
    (window as any).go = {
      main: { App: { ECDICTStatus: async () => ({ code: 500, err: 'download failed' }) } },
    };
    await expect(getECDICTStatus()).rejects.toThrow('download failed');
  });
});
