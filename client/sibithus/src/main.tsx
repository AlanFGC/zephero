import './index.css';

import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import App from './App.tsx';
import { AuthProvider } from './contexts/auth-context.tsx';

// biome-ignore lint/style/noNonNullAssertion: <we are sure that the root element exists>
createRoot(document.getElementById('root')!).render(
	<StrictMode>
		<AuthProvider>
		<App />
		</AuthProvider>
	</StrictMode>,
);
