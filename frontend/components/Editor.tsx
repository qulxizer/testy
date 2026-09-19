"use client";

import Editor from "@monaco-editor/react";

interface EditorPaneProps {
  files: Record<string, string>;
  activeTab: string;
  solutionFilename: string;
  onTabChange: (filename: string) => void;
  onCodeChange: (filename: string, code: string) => void;
}

export function EditorPane({
  files,
  activeTab,
  solutionFilename,
  onTabChange,
  onCodeChange,
}: EditorPaneProps) {
  const isApp = solutionFilename === "main.go";

  return (
    <div className="flex h-full w-full flex-col border-r border-zinc-800">
      {/* Tabs only show up for function quests */}
      {!isApp && (
        <div className="flex h-9 shrink-0 border-b border-zinc-800 bg-zinc-900/40">
          {Object.keys(files).map((fileName) => {
            const isActive = activeTab === fileName;
            return (
              <button
                key={fileName}
                onClick={() => onTabChange(fileName)}
                className={`flex items-center gap-2 border-r border-zinc-800 px-4 text-xs font-mono transition-colors ${
                  isActive
                    ? "border-b-2 border-b-emerald-500 bg-zinc-950 font-medium text-emerald-400"
                    : "text-zinc-400 hover:bg-zinc-800/30 hover:text-zinc-200"
                }`}
              >
                <span>{fileName}</span>
              </button>
            );
          })}
        </div>
      )}

      {/* Editor Container */}
      <div className="relative flex-1 min-h-0 w-full">
        <Editor
          height="100%"
          defaultLanguage="go"
          theme="vs-dark"
          path={activeTab}
          value={files[activeTab] || ""}
          onChange={(val) => onCodeChange(activeTab, val || "")}
          options={{
            minimap: { enabled: false },
            fontSize: 14,
            tabSize: 4,
            insertSpaces: false,
            automaticLayout: true, // Re-computes size smoothly on panel drag
          }}
        />
      </div>
    </div>
  );
}
