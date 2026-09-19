import Link from "next/link";

export default function HomePage() {
  return (
    <main className="min-h-screen bg-zinc-950 text-white flex flex-col justify-between selection:bg-zinc-800">
      {/* Navbar */}
      <header className="border-b border-zinc-900 px-6 py-4">
        <div className="mx-auto flex max-w-5xl items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-white font-bold text-black text-base">
              T
            </div>
            <span className="font-semibold tracking-tight text-zinc-100">
              Testy
            </span>
          </div>

          <Link
            href="/login"
            className="rounded-lg border border-zinc-800 bg-zinc-900 px-4 py-2 text-sm font-medium text-zinc-200 transition hover:border-zinc-700 hover:bg-zinc-800"
          >
            Sign in
          </Link>
        </div>
      </header>

      {/* Hero Section */}
      <section className="flex-1 flex flex-col items-center justify-center px-4 py-20 text-center">
        <div className="mx-auto max-w-2xl">
          <div className="mb-4 inline-flex items-center gap-2 rounded-full border border-zinc-800 bg-zinc-900/80 px-3 py-1 text-xs text-zinc-400">
            <span className="h-1.5 w-1.5 rounded-full bg-emerald-500 animate-pulse" />
            Synced with 01-edu / Reboot01
          </div>

          <h1 className="text-4xl font-bold tracking-tight sm:text-6xl text-white">
            Reboot01 code runner.
          </h1>

          <p className="mt-4 text-base text-zinc-400 max-w-md mx-auto">
            Just a harness wired up to run the official 01-edu / Reboot01 tests
            inside an isolated sandbox so you can practice for your checkpoints.
          </p>

          <div className="mt-8 flex items-center justify-center gap-3">
            <Link
              href="/login"
              className="rounded-lg bg-white px-5 py-2.5 text-sm font-medium text-black transition hover:bg-zinc-200"
            >
              Sign In
            </Link>

            <Link
              href="/testy"
              className="rounded-lg border border-zinc-800 bg-zinc-900 px-5 py-2.5 text-sm font-medium text-zinc-300 transition hover:border-zinc-700 hover:bg-zinc-800"
            >
              Open Workspace
            </Link>
          </div>
        </div>

        {/* Feature Cards */}
        <div className="mt-16 grid w-full max-w-3xl grid-cols-1 gap-4 px-4 sm:grid-cols-3 text-left">
          <div className="rounded-xl border border-zinc-800/80 bg-zinc-900/40 p-5">
            <div className="text-sm font-medium text-zinc-200">
              Official Tests
            </div>
            <p className="mt-2 text-xs text-zinc-400">
              Pulled straight from Reboot01 & 01-edu repos to match exact
              exercise specs.
            </p>
          </div>

          <div className="rounded-xl border border-zinc-800/80 bg-zinc-900/40 p-5">
            <div className="text-sm font-medium text-zinc-200">
              Kept Up to Date
            </div>
            <p className="mt-2 text-xs text-zinc-400">
              Tracked against curriculum updates so edge cases stay current.
            </p>
          </div>

          <div className="rounded-xl border border-zinc-800/80 bg-zinc-900/40 p-5">
            <div className="text-sm font-medium text-zinc-200">
              Isolated Runs
            </div>
            <p className="mt-2 text-xs text-zinc-400">
              Executes in a clean container and spits out stdout, stderr, and
              pass/fail diffs.
            </p>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-zinc-900 px-6 py-4 text-center text-xs text-zinc-500">
        Testy &bull; Unofficial test runner wired for Reboot01 / 01-edu
      </footer>
    </main>
  );
}
