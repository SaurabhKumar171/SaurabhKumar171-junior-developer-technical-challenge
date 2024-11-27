import React from 'react';
import './App.css';
import CharacterSearch from './components/core/CharacterSearch';

function App() {
  return (
    <div className="min-h-screen bg-gray-100">
      <header className="bg-blue-600 text-white p-4 text-center text-2xl">
        Rick and Morty Character Search
      </header>
      <main className="p-6">
        <CharacterSearch />
      </main>
    </div>

  );
}

export default App;
