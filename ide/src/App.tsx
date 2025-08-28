import { Button } from "../components/ui/button";
import Header from "./components/header";

function App() {
  return (
    <div className="w-screen h-screen flex flex-col bg-neutral-950 text-white">
      <Header />

      <main className="flex-1 p-6">
        <Button variant="outline">Hello World</Button>
      </main>
    </div>
  );
}

export default App;
