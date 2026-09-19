"use client";

import { ComponentPropsWithoutRef, useState } from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

interface SubjectProps {
  subject: string;
  audit?: string;
}

export function Subject({ subject, audit }: SubjectProps) {
  const [activeTab, setActiveTab] = useState<"subject" | "audit">("subject");

  const content = activeTab === "subject" ? subject : audit || "";

  return (
    <div className="flex h-full w-full flex-col border-r border-zinc-800 bg-zinc-950">
      {/* Header / Tabs */}
      <div className="flex h-9 shrink-0 border-b border-zinc-800 bg-zinc-900/40">
        <button
          onClick={() => setActiveTab("subject")}
          className={`flex items-center gap-2 border-r border-zinc-800 px-4 text-xs font-mono transition-colors ${
            activeTab === "subject"
              ? "border-b-2 border-b-emerald-500 bg-zinc-950 font-medium text-emerald-400"
              : "text-zinc-400 hover:bg-zinc-800/30 hover:text-zinc-200"
          }`}
        >
          <span>subject.md</span>
        </button>

        {audit && (
          <button
            onClick={() => setActiveTab("audit")}
            className={`flex items-center gap-2 border-r border-zinc-800 px-4 text-xs font-mono transition-colors ${
              activeTab === "audit"
                ? "border-b-2 border-b-emerald-500 bg-zinc-950 font-medium text-emerald-400"
                : "text-zinc-400 hover:bg-zinc-800/30 hover:text-zinc-200"
            }`}
          >
            <span>audit.md</span>
          </button>
        )}
      </div>

      {/* Markdown Body */}
      <div className="flex-1 overflow-y-auto p-6 text-sm text-zinc-300">
        <ReactMarkdown
          remarkPlugins={[remarkGfm]}
          components={{
            h1: ({ children, ...props }: ComponentPropsWithoutRef<"h1">) => (
              <h1
                className="mb-4 mt-2 border-b border-zinc-800 pb-2 text-xl font-bold text-zinc-100"
                {...props}
              >
                {children}
              </h1>
            ),
            h2: ({ children, ...props }: ComponentPropsWithoutRef<"h2">) => (
              <h2
                className="mb-3 mt-6 text-base font-semibold text-zinc-100"
                {...props}
              >
                {children}
              </h2>
            ),
            h3: ({ children, ...props }: ComponentPropsWithoutRef<"h3">) => (
              <h3
                className="mb-2 mt-4 text-sm font-semibold text-zinc-200"
                {...props}
              >
                {children}
              </h3>
            ),
            p: ({ children, ...props }: ComponentPropsWithoutRef<"p">) => (
              <p className="mb-4 leading-relaxed text-zinc-300" {...props}>
                {children}
              </p>
            ),
            ul: ({ children, ...props }: ComponentPropsWithoutRef<"ul">) => (
              <ul
                className="mb-4 list-disc space-y-1 pl-5 text-zinc-300"
                {...props}
              >
                {children}
              </ul>
            ),
            ol: ({ children, ...props }: ComponentPropsWithoutRef<"ol">) => (
              <ol
                className="mb-4 list-decimal space-y-1 pl-5 text-zinc-300"
                {...props}
              >
                {children}
              </ol>
            ),
            pre: ({ children, ...props }: ComponentPropsWithoutRef<"pre">) => (
              <pre
                className="my-4 overflow-x-auto rounded-md border border-zinc-800 bg-zinc-900/60 p-3 font-mono text-xs leading-relaxed text-zinc-200"
                {...props}
              >
                {children}
              </pre>
            ),
            code: ({
              children,
              className,
              ...props
            }: ComponentPropsWithoutRef<"code">) => {
              const isInline =
                !className &&
                typeof children === "string" &&
                !children.includes("\n");

              if (isInline) {
                return (
                  <code
                    className="rounded bg-zinc-800/80 px-1.5 py-0.5 font-mono text-xs text-amber-300"
                    {...props}
                  >
                    {children}
                  </code>
                );
              }

              return (
                <code
                  className={`block font-mono text-xs text-zinc-200 ${className || ""}`}
                  {...props}
                >
                  {children}
                </code>
              );
            },
            table: ({
              children,
              ...props
            }: ComponentPropsWithoutRef<"table">) => (
              <div className="my-4 overflow-x-auto">
                <table
                  className="w-full border-collapse border border-zinc-800 text-xs"
                  {...props}
                >
                  {children}
                </table>
              </div>
            ),
            th: ({ children, ...props }: ComponentPropsWithoutRef<"th">) => (
              <th
                className="border border-zinc-800 bg-zinc-900 px-3 py-2 text-left font-semibold text-zinc-200"
                {...props}
              >
                {children}
              </th>
            ),
            td: ({ children, ...props }: ComponentPropsWithoutRef<"td">) => (
              <td
                className="border border-zinc-800 px-3 py-2 text-zinc-300"
                {...props}
              >
                {children}
              </td>
            ),
          }}
        >
          {content}
        </ReactMarkdown>
      </div>
    </div>
  );
}
