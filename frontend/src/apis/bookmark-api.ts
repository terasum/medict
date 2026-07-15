// Bookmark & notebook API wrappers (#643) — typed IPC via wailsjs bindings.
import * as App from '../../wailsjs/go/main/App';
import { model } from './model';

export interface Bookmark {
    word: string;
    dict_id: string;
    dict_name: string;
    notebook_id: string;
    saved_at: number;
}

export interface Notebook {
    id: string;
    name: string;
    is_default: boolean;
    created_at: number;
}

// unwrap turns a backend error response (code != 200) into a thrown Error so
// callers can surface it via try/catch + message. Success responses yield data.
const unwrap = async <T>(p: Promise<model.Resp>): Promise<T> => {
    const resp = await p;
    if (resp.code !== 200) {
        throw new Error(resp.err || '操作失败');
    }
    return resp.data as T;
};

export const getBookmarks = async (): Promise<Bookmark[]> => {
    const resp = await App.GetBookmarks();
    return (resp.data as Bookmark[]) || [];
};

export const addBookmark = async (
    word: string,
    dictId: string,
    notebookId: string
): Promise<void> => {
    await unwrap<void>(App.AddBookmark(word, dictId, notebookId));
};

export const removeBookmark = async (word: string, dictId: string, notebookId: string): Promise<void> => {
    await unwrap<void>(App.RemoveBookmark(word, dictId, notebookId));
};

// getBookmarkSnapshot returns the stored self-contained HTML for a bookmark
// ("" if the bookmark has no snapshot — e.g. created before snapshots existed).
export const getBookmarkSnapshot = async (
    word: string,
    dictId: string,
    notebookId: string
): Promise<string> => {
    return unwrap<string>(App.GetBookmarkSnapshot(word, dictId, notebookId));
};

export const getNotebooks = async (): Promise<Notebook[]> => {
    const resp = await App.GetNotebooks();
    return (resp.data as Notebook[]) || [];
};

export const createNotebook = async (name: string): Promise<Notebook> => {
    return unwrap<Notebook>(App.CreateNotebook(name));
};

export const renameNotebook = async (id: string, name: string): Promise<void> => {
    await unwrap<void>(App.RenameNotebook(id, name));
};

export const deleteNotebook = async (id: string): Promise<void> => {
    await unwrap<void>(App.DeleteNotebook(id));
};

export const setDefaultNotebook = async (id: string): Promise<void> => {
    await unwrap<void>(App.SetDefaultNotebook(id));
};

// exportAnkiToApkg builds a native Anki .apkg from the given notebook (all
// notebooks if notebookId is empty) and writes it to a user-chosen path.
// Returns the saved path, or '' if the user cancelled the save dialog.
// Throws on backend error (code != 200).
export const exportAnkiToApkg = async (notebookId: string): Promise<string> => {
    const resp = await App.ExportAnki(notebookId);
    if (resp.code !== 200) {
        throw new Error(resp.err || '导出失败');
    }
    return (resp.data as string) || '';
};
