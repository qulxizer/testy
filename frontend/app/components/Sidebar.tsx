'use client';

import { useState, useMemo } from 'react';
import { Exercise } from '@/types';

interface SidebarProps {
  exercises: Exercise[];
  selectedExerciseKey: string;
  onSelectExercise: (ex: Exercise) => void;
}

export function Sidebar({ exercises, selectedExerciseKey, onSelectExercise }: SidebarProps) {
  const [search, setSearch] = useState('');
  const [open, setOpen] = useState<Record<string, boolean>>({});

  const toggle = (key: string) => {
    setOpen((prev) => ({ ...prev, [key]: !prev[key] }));
  };

  const { quests, checkpoints } = useMemo(() => {
    const filtered = exercises.filter(
      (ex) =>
        ex.name.toLowerCase().includes(search.toLowerCase()) ||
        ex.key.toLowerCase().includes(search.toLowerCase())
    );

    const qMap: Record<string, Exercise[]> = {};
    const cpMap: Record<string, Record<string, Exercise[]>> = {};

    filtered.forEach((ex) => {
      const parts = (ex.path || '').split('/').filter(Boolean);
      const cpIdx = parts.findIndex((p) => p.includes('checkpoint'));
      const qIdx = parts.findIndex((p) => p.includes('quest'));

      if (cpIdx !== -1) {
        const cpName = parts[cpIdx];
        const diffKey = `Difficulty ${ex.difficulty ? ex.difficulty / 2 : 0}`;

        if (!cpMap[cpName]) cpMap[cpName] = {};
        if (!cpMap[cpName][diffKey]) cpMap[cpName][diffKey] = [];
        cpMap[cpName][diffKey].push(ex);
      } else {
        const qName = qIdx !== -1 ? parts[qIdx] : 'other';
        if (!qMap[qName]) qMap[qName] = [];
        qMap[qName].push(ex);
      }
    });

    return { quests: qMap, checkpoints: cpMap };
  }, [exercises, search]);

  return (
    <aside className="flex w-64 flex-col border-r border-zinc-800 bg-zinc-900/50">
      <div className="p-3 border-b border-zinc-800">
        <input
          type="text"
          placeholder="Search exercises..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="w-full rounded bg-zinc-800 px-2.5 py-1 text-xs text-zinc-200 placeholder-zinc-500 outline-none focus:ring-1 focus:ring-zinc-600"
        />
      </div>

      <div className="flex-1 overflow-y-auto font-mono text-xs select-none">
        {/* QUESTS */}
        {Object.keys(quests).length > 0 && (
          <div className="mb-4">
            <div className="px-3 py-1.5 text-[10px] font-bold uppercase tracking-wider text-zinc-500 bg-zinc-950/60 sticky top-0 border-b border-zinc-800/40">
              Quests
            </div>
            {Object.keys(quests).sort().map((qKey) => {
              const items = quests[qKey];
              const isOpen = Boolean(open[qKey]);
              return (
                <div key={qKey} className="border-b border-zinc-800/30">
                  <button
                    onClick={() => toggle(qKey)}
                    className="flex w-full items-center justify-between px-3 py-1.5 text-zinc-300 hover:bg-zinc-800/50 font-semibold"
                  >
                    <span className="flex items-center gap-1.5">
                      <span className="text-[9px] text-zinc-500">{isOpen ? '▼' : '▶'}</span>
                      {qKey}
                    </span>
                    <span className="text-[10px] text-zinc-500">{items.length}</span>
                  </button>

                  {isOpen && (
                    <div className="bg-zinc-950/30">
                      {items.map((ex) => (
                        <button
                          key={ex.id || ex.key}
                          onClick={() => onSelectExercise(ex)}
                          className={`w-full text-left pl-6 pr-3 py-1.5 transition-colors ${
                            ex.key === selectedExerciseKey
                              ? 'bg-zinc-800 font-semibold text-emerald-400'
                              : 'text-zinc-400 hover:bg-zinc-800/40 hover:text-zinc-200'
                          }`}
                        >
                          {ex.name}
                        </button>
                      ))}
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        )}

        {/* CHECKPOINTS */}
        {Object.keys(checkpoints).length > 0 && (
          <div className="mb-4">
            <div className="px-3 py-1.5 text-[10px] font-bold uppercase tracking-wider text-zinc-500 bg-zinc-950/60 sticky top-0 border-b border-zinc-800/40">
              Checkpoints
            </div>
            {Object.keys(checkpoints).sort().map((cpKey) => {
              const levels = checkpoints[cpKey];
              const isCpOpen = Boolean(open[cpKey]);

              return (
                <div key={cpKey} className="border-b border-zinc-800/30">
                  <button
                    onClick={() => toggle(cpKey)}
                    className="flex w-full items-center justify-between px-3 py-1.5 text-zinc-200 hover:bg-zinc-800/50 font-semibold"
                  >
                    <span className="flex items-center gap-1.5">
                      <span className="text-[9px] text-zinc-500">{isCpOpen ? '▼' : '▶'}</span>
                      {cpKey}
                    </span>
                  </button>

                  {isCpOpen && (
                    <div className="bg-zinc-950/20">
                      {Object.keys(levels).sort().map((lvlKey) => {
                        const items = levels[lvlKey];
                        const lvlPathKey = `${cpKey}/${lvlKey}`;
                        const isLvlOpen = Boolean(open[lvlPathKey]);

                        return (
                          <div key={lvlKey}>
                            <button
                              onClick={() => toggle(lvlPathKey)}
                              className="flex w-full items-center justify-between pl-5 pr-3 py-1 text-zinc-400 hover:bg-zinc-800/40 hover:text-zinc-200"
                            >
                              <span className="flex items-center gap-1.5">
                                <span className="text-[8px] text-zinc-600">
                                  {isLvlOpen ? '▼' : '▶'}
                                </span>
                                <span className="font-medium text-[11px]">{lvlKey}</span>
                              </span>
                              <span className="text-[9px] text-zinc-600">{items.length}</span>
                            </button>

                            {isLvlOpen && (
                              <div className="bg-zinc-950/40">
                                {items.map((ex) => (
                                  <button
                                    key={ex.id || ex.key}
                                    onClick={() => onSelectExercise(ex)}
                                    className={`w-full text-left pl-8 pr-3 py-1 transition-colors ${
                                      ex.key === selectedExerciseKey
                                        ? 'bg-zinc-800 font-semibold text-emerald-400'
                                        : 'text-zinc-500 hover:bg-zinc-800/30 hover:text-zinc-300'
                                    }`}
                                  >
                                    {ex.name}
                                  </button>
                                ))}
                              </div>
                            )}
                          </div>
                        );
                      })}
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        )}
      </div>
    </aside>
  );
}
