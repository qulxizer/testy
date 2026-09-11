interface OutputPaneProps {
  output: string;
}

export function OutputPane({ output }: OutputPaneProps) {
  return (
    <div className="h-full w-80 overflow-auto bg-black p-4 font-mono text-xs">
      <div className="text-zinc-500">// Terminal Output</div>
      <pre className="mt-2 whitespace-pre-wrap text-zinc-200">{output}</pre>
    </div>
  );
}
