interface HeaderProps {
  exerciseName: string;
  filename: string;
  loading: boolean;
  onRunTests: () => void;
}

export function Header({ exerciseName, filename, loading, onRunTests }: HeaderProps) {
  return (
    <header className="flex h-12 items-center justify-between border-b border-zinc-800 px-4">
      <div className="flex items-center gap-3">
        <span className="font-mono text-sm font-semibold">{exerciseName}</span>
        <span className="rounded bg-zinc-800 px-2 py-0.5 font-mono text-xs text-zinc-400">
          {filename}
        </span>
      </div>
      <button
        onClick={onRunTests}
        disabled={loading}
        className="rounded bg-emerald-600 px-4 py-1.5 text-sm font-medium hover:bg-emerald-500 disabled:opacity-50"
      >
        {loading ? 'Testing...' : 'Submit & Test'}
      </button>
    </header>
  );
}
