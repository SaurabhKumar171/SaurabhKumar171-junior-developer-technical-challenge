import React, { useState } from 'react';
import axios from 'axios';
import Spinner from '../../common/Spinner';

// Define types for character and episode data
interface Origin {
  name: string;
  url: string;
}

interface Location {
  name: string;
  url: string;
}

interface Character {
  id: number;
  name: string;
  status: string;
  species: string;
  gender: string;
  origin: Origin;
  location: Location;
  image: string;
  episode: string[];
  url: string;
  created: string;
}

const CharacterSearch = () => {
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [characters, setCharacters] = useState<Character[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>('');
  const [showAll, setShowAll] = useState<boolean>(false);

  // Pagination states
  const [page, setPage] = useState<number>(1);
  const [limit, setLimit] = useState<number>(6);
  const [totalPages, setTotalPages] = useState<number>(1);

  const handleSearch = async () => {
    setCharacters([]);
    if (!searchQuery) return;
    setLoading(true);
    setError('');
    
    try {
      const response = await axios.get(`${process.env.REACT_APP_API_URL}/characters/search?name=${searchQuery}`);
      setCharacters(response.data.characters);
    } catch (err) {
      setError('Error fetching characters');
    } finally {
      setLoading(false);
    }
  };

  const handleShowAll = async () => {
    setCharacters([]);
    setLoading(true);
    setError('');
    setShowAll(true);

    try {
      const response = await axios.get(`${process.env.REACT_APP_API_URL}/getCharacters?page=${page}&limit=${limit}`);  
      setCharacters(response.data.characters);
      setTotalPages(Math.ceil(response.data.total / limit));
    } catch (err) {
      setError('Error fetching all characters');
    } finally {
      setLoading(false);
    }
  };

  const handlePageChange = (newPage: number) => {
    if (newPage < 1 || newPage > totalPages) return;
    setPage(newPage);
    handleShowAll();
  };

  return (
    <div className="max-w-5xl mx-auto p-6 relative">
      <div className="flex mb-4 justify-between items-center">
        <div className="flex space-x-4">
          <input
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            placeholder="Search for a character"
            className="p-2 w-full border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <button
            onClick={handleSearch}
            className="ml-2 p-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition"
          >
            Search
          </button>
        </div>
        <button
          onClick={handleShowAll}
          className="p-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition"
        >
          Show All Characters
        </button>
      </div>

      {loading && <Spinner />}
      {error && <p className="text-center text-red-500">{error}</p>}

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 mt-6">
        {characters.length > 0 ? (
          characters.map((character) => (
            <div
              key={character.id}
              className="border border-gray-300 p-4 rounded-lg shadow-lg hover:shadow-xl transition transform hover:scale-105 h-auto"
            >
              <img
                src={character.image}
                alt={character.name}
                className="w-full h-48 object-cover rounded-md mb-4"
              />
              <h3 className="text-xl font-semibold text-gray-800">{character.name}</h3>
              <p className="text-gray-600">Status: {character.status}</p>
              <p className="text-gray-600">Species: {character.species}</p>
              <p className="text-gray-600">Gender: {character.gender}</p>

              <div className="mt-2">
                <p className="text-sm text-gray-500">Origin: {character.origin.name}</p>
                <p className="text-sm text-gray-500">Location: {character.location.name}</p>
              </div>

              <h4 className="mt-4 font-semibold text-gray-800">Episodes:</h4>
              <div className="flex flex-wrap gap-2 text-sm text-gray-600 overflow-y-scroll max-h-20">
                {character.episode.map((episode, index) => (
                  <span key={index}>
                    <a
                      href={episode}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-blue-600 hover:underline"
                    >
                      Episode {index + 1}
                    </a>
                  </span>
                ))}
              </div>
            </div>
          ))
        ) : (
          !loading &&
          !error && (
            <p className="text-center text-gray-500 mt-6">
              No characters found. Please try searching with a different keyword.
            </p>
          )
        )}
      </div>


      { characters.length > 0 &&
        <div className="flex justify-center mt-6 space-x-4">
        <button
          onClick={() => handlePageChange(page - 1)}
          disabled={page === 1}
          className="p-2 bg-gray-600 text-white rounded-md hover:bg-gray-700 transition"
        >
          Previous
        </button>
        <span className="self-center text-gray-700">Page {page} of {totalPages}</span>
        <button
          onClick={() => handlePageChange(page + 1)}
          disabled={page === totalPages}
          className="p-2 bg-gray-600 text-white rounded-md hover:bg-gray-700 transition"
        >
          Next
        </button>
      </div>}
    </div>
  );
};

export default CharacterSearch;
