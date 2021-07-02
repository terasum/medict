import { LinkReplacer } from '../domain/ReplacerLink';
import { ImageReplacer } from '../domain/ReplacerImage';
import { SoundReplacer } from '../domain/ReplacerSound';
import { CSSReplacer } from '../domain/ReplacerCSS';
import { EntryReplacer } from '../domain/ReplacerEntry';
import { JSReplacer } from '../domain/ReplacerJS';

import { ResourceFn, LookupFn } from '../domain/Replacer';

const replacerChain = [
  new LinkReplacer(),
  new ImageReplacer(),
  new SoundReplacer(),
  new CSSReplacer(),
  new JSReplacer(),
  new EntryReplacer(),
];

export const dictContentService = {
  definitionReplace: (
    dictid: string,
    keyText: string,
    html: string,
    lookupFn: LookupFn,
    resourceFn: ResourceFn
  ) => {
    replacerChain.forEach(replacer => {
      html = replacer.replace(dictid, keyText, html, lookupFn, resourceFn);
    });
    return html;
  },
};
