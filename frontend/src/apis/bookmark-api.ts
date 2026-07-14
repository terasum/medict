// Bookmark API wrappers (#643) — typed IPC via wailsjs bindings.
import * as App from '../../wailsjs/go/main/App';

export interface Bookmark {
    word: string;
    dict_id: string;
    dict_name: string;
    saved_at: number;
}

export const getBookmarks = async (): Promise<Bookmark[]> => {
    const resp = await App.GetBookmarks();
    return (resp.data as Bookmark[]) || [];
};

export const addBookmark = async (word: string, dictId: string, dictName: string): Promise<void> => {
    await App.AddBookmark(word, dictId, dictName);
};

export const removeBookmark = async (word: string, dictId: string): Promise<void> => {
    await App.RemoveBookmark(word, dictId);
};
