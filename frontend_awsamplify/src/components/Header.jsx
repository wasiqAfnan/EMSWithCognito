export default function Header({ userEmail, onSignOut }) {
  return (
    <header className="bg-white shadow">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
        <h1 className="text-2xl font-bold">EMS Dashboard</h1>
        <div className="flex items-center gap-4">
          <span className="text-sm text-gray-600">Hi, {userEmail}</span>
          <button
            onClick={onSignOut}
            className="px-4 py-2 rounded-md shadow-sm text-sm font-medium text-white bg-red-600 hover:bg-red-700 focus:outline-none cursor-pointer"
          >
            Sign out
          </button>
        </div>
      </div>
    </header>
  );
}

