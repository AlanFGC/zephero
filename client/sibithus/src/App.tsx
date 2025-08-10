import './App.css';
import { AuthCard } from './components/auth-card';
import { useAuth } from './contexts/auth-context';

const App = () => {
	const authClient = useAuth();
	const session = authClient.useSession();
	
	return (
		<div className="flex min-h-svh flex-col items-center justify-center">
			{!session.data?.session &&
				<AuthCard />}
		</div>
	);
}

export default App;
