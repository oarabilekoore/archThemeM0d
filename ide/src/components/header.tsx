import { Bolt } from "lucide-react";

export default function Header() {
  return (
    <header className="w-full bg-neutral-900 text-white shadow-md">
      <div className="flex items-center justify-between px-6 py-3">
        <a className="flex items-center gap-3 cursor-pointer">
          <img
            src="/archThemeM0d.png"
            alt="archThemeM0d logo"
            className="h-8 w-8"
          />
          <span className="text-lg font-semibold">archThemeM0dIDE</span>
        </a>

        <div className="flex items-center gap-4">
          <Bolt className="w-5 h-5 cursor-pointer" />
        </div>
      </div>
    </header>
  );
}
