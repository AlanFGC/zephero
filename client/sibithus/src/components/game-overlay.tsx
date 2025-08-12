import type React from 'react';
import type { ReactNode } from 'react';
import { useAuth } from '../contexts/auth-context';

interface GameOverlayProps {
  children?: ReactNode;
}

const GameOverlay: React.FC<GameOverlayProps> = ({ children }) => {
  const auth = useAuth();

  const handleSignOut = () => {
    console.log('You have signed out successfully.');
    auth.signOut();
  };

  return (
    <div>
      <div className="absolute bottom-0 left-1/2 transform -translate-x-1/2 bg-red-500 rounded-t-lg w-[300px] h-[5vh] z-10 flex items-center justify-center">
        <button 
          type="button"
          onClick={handleSignOut}
          className="text-white font-semibold hover:text-gray-200 transition-colors"
        >
          Sign Out
        </button>
      </div>
      {children}
    </div>
  );
};

export {
  GameOverlay
};