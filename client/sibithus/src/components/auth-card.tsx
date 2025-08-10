import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { Button } from './ui/button';
import { LoginForm } from './login-form';
import { SignUpForm } from './sign-up-form';

export function AuthCard() {
  const [isLogin, setIsLogin] = useState(true);

  return (
    <Card className="w-full max-w-sm">
      <CardHeader className="relative">
        <CardTitle>{isLogin ? 'Login to your account' : 'Create an account'}</CardTitle>
        <Button
          variant="ghost"
          size="sm"
          className="absolute top-4 right-4"
          onClick={() => setIsLogin(!isLogin)}
        >
          {isLogin ? 'Sign Up' : 'Login'}
        </Button>
      </CardHeader>
      <CardContent>
        {isLogin ? <LoginForm /> : <SignUpForm />}
      </CardContent>
    </Card>
  );
}