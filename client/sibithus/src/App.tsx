import './App.css';
import { AuthCard } from './components/auth-card';
import { GameOverlay } from './components/game-overlay';
import { useAuth } from './contexts/auth-context';

const App = () => {
	const authClient = useAuth();
	const { data: session, isPending } = authClient.useSession();
	
	if (isPending) {
		return (
			<div className="flex min-h-svh flex-col items-center justify-center">
				<div>Loading...</div>
			</div>
		);
	}
	
	return (
		<div className="flex min-h-svh flex-col items-center justify-center">
			{!session &&
				<AuthCard />}
			{session &&
				<GameOverlay>
					<canvas id="gameCanvas" className="w-full h-full"></canvas>
				</GameOverlay>}
		</div>
	);
}

export default App;
