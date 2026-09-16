"use client";

import { useState } from "react";
import exercisesData from "@/data/go_exercises.json";
import { Exercise } from "@/lib/types";
import { getInitialFiles } from "@/lib/editor-utils";
import { Header } from "@/components/Header";
import { Sidebar } from "@/components/Sidebar";
import { EditorPane } from "@/components/Editor";
import { OutputPane } from "@/components/Output";

const exercises: Exercise[] = exercisesData as Exercise[];

const defaultExercise: Exercise = exercises[0] || {
  id: 1,
  key: "onlya",
  name: "onlya",
  path: "/bahrain/bh-piscine/quest-01/onlya",
  language: "go",
  expectedFiles: ["main.go"],
};

export default function WorkspacePage() {
  const [selectedExercise, setSelectedExercise] =
    useState<Exercise>(defaultExercise);
  const [files, setFiles] = useState<Record<string, string>>(() =>
    getInitialFiles(defaultExercise),
  );
  const [activeTab, setActiveTab] = useState<string>(
    defaultExercise.expectedFiles?.[0] || "main.go",
  );
  const [output, setOutput] = useState("");
  const [loading, setLoading] = useState(false);

  const solutionFilename = selectedExercise.expectedFiles?.[0] || "main.go";

  const handleSelectExercise = (ex: Exercise) => {
    const initFiles = getInitialFiles(ex);
    setSelectedExercise(ex);
    setFiles(initFiles);
    setActiveTab(ex.expectedFiles?.[0] || "main.go");
    setOutput("");
  };

  const handleCodeChange = (fileName: string, value: string) => {
    setFiles((prev) => ({ ...prev, [fileName]: value }));
  };

  const runTests = async () => {
    setLoading(true);
    setOutput("Running test suite...");

    try {
      const res = await fetch("/api/test", {
        credentials: "include",
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          exercise: selectedExercise.key,
          filename: solutionFilename,
          code: files[solutionFilename],
          main: files["main.go"] ?? null,
          testImage: selectedExercise.testImage,
        }),
      });

      const text = await res.text();
      try {
        const data = JSON.parse(text);
        setOutput(data.output || "No output.");
      } catch {
        setOutput(`Server error (${res.status}):\n${text}`);
      }
    } catch (err: any) {
      setOutput(`Failed to reach test runner: ${err?.message || err}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex h-screen flex-col bg-zinc-950 text-white">
      <Header
        exerciseName={selectedExercise.name}
        filename={solutionFilename}
        loading={loading}
        onRunTests={runTests}
      />
      <div className="flex flex-1 overflow-hidden">
        <Sidebar
          exercises={exercises}
          selectedExerciseKey={selectedExercise.key}
          onSelectExercise={handleSelectExercise}
        />
        <EditorPane
          files={files}
          activeTab={activeTab}
          solutionFilename={solutionFilename}
          onTabChange={setActiveTab}
          onCodeChange={handleCodeChange}
        />
        <OutputPane output={output} />
      </div>
    </div>
  );
}
