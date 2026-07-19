interface BackendResponse {
  code: number;
  data?: unknown;
  err?: string;
}

export interface ECDICTStatus {
  entryCount: number;
  edition: 'compact' | 'full';
}

function appMethod(name: 'ECDICTStatus' | 'InstallFullECDICT'): () => Promise<BackendResponse> {
  const method = (window as any)?.go?.main?.App?.[name];
  if (typeof method !== 'function') {
    throw new Error('完整词库功能仅在 Medict 桌面应用中可用');
  }
  return method;
}

function unwrapStatus(response: BackendResponse): ECDICTStatus {
  if (response.code !== 200) {
    throw new Error(response.err || 'ECDICT 操作失败');
  }
  return response.data as ECDICTStatus;
}

export async function getECDICTStatus(): Promise<ECDICTStatus> {
  return unwrapStatus(await appMethod('ECDICTStatus')());
}

export async function installFullECDICT(): Promise<ECDICTStatus> {
  return unwrapStatus(await appMethod('InstallFullECDICT')());
}
