/**
 *
 * Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

export type OSPlatform = 'darwin' | 'windows' | 'linux' | 'unknown';

// Synchronous user-agent sniff. Good enough to paint the correct layout on
// first render; detectPlatform() refines it with the authoritative answer
// from the Go side.
export function sniffPlatformSync(): OSPlatform {
  const ua: string = navigator.userAgent ?? '';
  if (/Windows/.test(ua)) return 'windows';
  if (/Macintosh|Mac OS X/.test(ua)) return 'darwin';
  if (/Linux|X11/.test(ua)) return 'linux';
  return 'unknown';
}

// Resolve the host OS via the App.Platform() Wails binding (runtime.GOOS),
// falling back to the user agent when the bridge is unavailable (e.g. the
// plain vite dev server in a browser).
export async function detectPlatform(): Promise<OSPlatform> {
  const call = (window as any)?.go?.main?.App?.Platform;
  if (typeof call === 'function') {
    try {
      const goos = await call();
      if (goos === 'darwin' || goos === 'windows' || goos === 'linux') {
        return goos;
      }
    } catch {
      // bridge hiccup: fall through to the user-agent guess
    }
  }
  return sniffPlatformSync();
}

// Class bound on #app-root; the SCSS uses the .os-windows / .os-linux
// variants to shorten the fake title bar where the native one cannot be
// hidden.
export function platformClass(platform: OSPlatform): string {
  return `os-${platform}`;
}
