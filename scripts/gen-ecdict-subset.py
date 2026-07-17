#!/usr/bin/env python3
# Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
# GPL-3.0
#
# Generate a compact ECDICT subset (top-N frequent English->Chinese entries) as a
# SQLite file for embedding as Medict's default offline EN-CN dictionary (#782/roadmap).
#
# Source: https://github.com/skywind3000/ECDICT  (ecdict.csv, CC-BY-SA / public)
# We keep only: word, phonetic, definition (EN), translation (CN), pos, frq — and
# only rows that HAVE a Chinese translation — then take the top N by frequency
# (collins stars first, then frq). Output ~8MB for N=50000.
#
# Usage: python3 gen-ecdict-subset.py [ecdict.csv] [out.db] [N]
import csv
import os
import sqlite3
import sys

SRC = sys.argv[1] if len(sys.argv) > 1 else "/tmp/ecdict.csv"
OUT = sys.argv[2] if len(sys.argv) > 2 else "internal/entry/preset/ecdict/ecdict.db"
N = int(sys.argv[3]) if len(sys.argv) > 3 else 50000


def to_int(v):
    try:
        return int(v or 0)
    except ValueError:
        return 0


def main():
    rows = []
    with open(SRC, newline="", encoding="utf-8") as f:
        for r in csv.DictReader(f):
            word = (r.get("word") or "").strip()
            trans = (r.get("translation") or "").strip()
            if not word or not trans:  # need headword + Chinese translation
                continue
            collins = to_int(r.get("collins"))
            frq = to_int(r.get("frq"))
            rows.append((collins, frq, word, r.get("phonetic") or "",
                         r.get("definition") or "", trans, r.get("pos") or ""))

    # collins stars (1-5, higher=common) dominate; frq is a BNC frequency RANK
    # (lower=more frequent), so within a collins tier sort frq ascending.
    rows.sort(key=lambda x: (-x[0], x[1]))
    top = rows[:N]

    os.makedirs(os.path.dirname(OUT), exist_ok=True)
    if os.path.exists(OUT):
        os.remove(OUT)
    con = sqlite3.connect(OUT)
    con.execute(
        "CREATE TABLE ecdict (word TEXT PRIMARY KEY, phonetic TEXT, "
        "definition TEXT, translation TEXT, pos TEXT, frq INTEGER)"
    )
    con.executemany(
        "INSERT INTO ecdict(word,phonetic,definition,translation,pos,frq) VALUES(?,?,?,?,?,?)",
        [(w, ph, de, tr, po, fr) for (_cl, fr, w, ph, de, tr, po) in top],
    )
    con.commit()
    con.close()
    print(f"wrote {len(top)} rows to {OUT} ({os.path.getsize(OUT) // 1024} KB)")


if __name__ == "__main__":
    main()
