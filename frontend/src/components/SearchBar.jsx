export default function SearchBar({ searchQuery, setSearchQuery, onSearch, onClear, onAddEmployee }) {
  return (
    <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 bg-white p-4 rounded-lg shadow-sm">
      <form onSubmit={onSearch} className="flex gap-2 w-full sm:w-auto">
        <input
          type="text"
          placeholder="Search employees..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 min-w-[250px]"
        />
        <button
          type="submit"
          className="px-4 py-2 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200 font-medium border border-gray-300 cursor-pointer"
        >
          Search
        </button>
        {searchQuery && (
          <button
            type="button"
            onClick={onClear}
            className="px-3 py-2 text-gray-500 hover:text-gray-700 cursor-pointer"
          >
            Clear
          </button>
        )}
      </form>

      <button
        onClick={onAddEmployee}
        className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 font-medium shadow-sm cursor-pointer whitespace-nowrap"
      >
        + Add Employee
      </button>
    </div>
  );
}

