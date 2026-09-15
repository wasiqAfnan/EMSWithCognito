export default function Toast({ toast, onClose }) {
  if (!toast.visible) return null;

  return (
    <div className={`fixed top-4 left-1/2 transform -translate-x-1/2 z-50 px-6 py-3 rounded shadow-lg flex items-center gap-4 text-white font-medium transition-all duration-300 ${toast.type === 'error' ? 'bg-red-600' : 'bg-green-600'}`}>
      <span>{toast.message}</span>
      <button 
        onClick={onClose} 
        className="text-white hover:text-gray-200 cursor-pointer font-bold"
      >
        &times;
      </button>
    </div>
  );
}

