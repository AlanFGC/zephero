import { useState } from "react";
import { useDebouncedCallback } from "use-debounce";
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useAuth } from "@/contexts/auth-context";

type ValidationStatus = 'idle' | 'validating' | 'available' | 'taken';
type PasswordValidation = 'idle' | 'weak' | 'strong' | 'no_match';

const SignUpForm = () => {
  const authClient = useAuth();
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [usernameStatus, setUsernameStatus] = useState<ValidationStatus>('idle');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [passwordStatus, setPasswordStatus] = useState<PasswordValidation>('idle');
  const [confirmPasswordStatus, setConfirmPasswordStatus] = useState<PasswordValidation>('idle');
  
  const validatePassword = (password: string) => {
    const hasLetters = /[a-zA-Z]/.test(password);
    const hasNumbers = /\d/.test(password);
    
    if (!password.trim()) {
      return 'idle';
    }
    
    if (hasLetters && hasNumbers && password.length >= 6) {
      return 'strong';
    } else {
      return 'weak';
    }
  };

  const validatePasswordMatch = (password: string, confirmPassword: string) => {
    if (!confirmPassword.trim()) {
      return 'idle';
    }
    
    return password === confirmPassword ? 'strong' : 'no_match';
  };

  const handlePasswordChange = (value: string) => {
    setPassword(value);
    const validation = validatePassword(value) as PasswordValidation;
    setPasswordStatus(validation);
    
    if (confirmPassword) {
      const matchValidation = validatePasswordMatch(value, confirmPassword) as PasswordValidation;
      setConfirmPasswordStatus(matchValidation);
    }
  };

  const handleConfirmPasswordChange = (value: string) => {
    setConfirmPassword(value);
    const matchValidation = validatePasswordMatch(password, value) as PasswordValidation;
    setConfirmPasswordStatus(matchValidation);
  };

  const checkUsernameDebounced = useDebouncedCallback(
    async (value) => {
      if (!value.trim()) {
        setUsernameStatus('idle');
        return;
      }
      
      setUsernameStatus('validating');
      
      const response = await authClient.isUsernameAvailable({
      username: value})
      
      if (response?.error) {
        console.error("Error checking username availability:", response);
        setUsernameStatus('idle');
        return;
      }
      
      if (!response?.data?.available) {
        setUsernameStatus('taken');
      } else if (response?.data?.available) {
        setUsernameStatus('available');
      }
    },
    1000,
  )

  return (
    <form
      onSubmit={async (e) => {
        e.preventDefault();
        
        if (usernameStatus !== 'available') {
          alert('Please choose an available username.');
          return;
        }
        
        if (passwordStatus !== 'strong') {
          alert('Please enter a strong password.');
          return;
        }
        
        if (confirmPasswordStatus !== 'strong') {
          alert('Passwords do not match.');
          return;
        }
        

        await authClient.signUp.email({
          email,
          password,
          name: username,
          username: username,
      });


        
      }}
    >
      <div className="flex flex-col gap-6">
        <div className="grid gap-2">
          <Label htmlFor="username">Username</Label>
          {usernameStatus === 'validating' && (
            <p className="text-sm text-gray-500">Checking availability...</p>
          )}
          {usernameStatus === 'available' && (
            <p className="text-sm text-green-600">Username is available</p>
          )}
          {usernameStatus === 'taken' && (
            <p className="text-sm text-red-500">Username is already taken</p>
          )}
          <Input
            id="username"
            type="text"
            placeholder="Enter username"
            required
            value={username}
            onChange={(e) => {
              const value = e.target.value;
              setUsername(value);
              checkUsernameDebounced(value);
            }}
          />
        </div>
        <div className="grid gap-2">
          <Label htmlFor="email">Email</Label>
          <Input
            id="email"
            type="email"
            placeholder="m@example.com"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
        </div>
        <div className="grid gap-2">
          <Label htmlFor="password">Password</Label>
          {passwordStatus === 'weak' && (
            <p className="text-sm text-red-500">Password must contain letters, numbers, and be at least 6 characters</p>
          )}
          {passwordStatus === 'strong' && (
            <p className="text-sm text-green-600">Password is strong</p>
          )}
          <Input 
            id="password" 
            type="password" 
            required
            value={password}
            onChange={(e) => handlePasswordChange(e.target.value)}
          />
        </div>
        <div className="grid gap-2">
          <Label htmlFor="confirmPassword">Confirm Password</Label>
          {confirmPasswordStatus === 'no_match' && (
            <p className="text-sm text-red-500">Passwords do not match</p>
          )}
          {confirmPasswordStatus === 'strong' && (
            <p className="text-sm text-green-600">Passwords match</p>
          )}
          <Input 
            id="confirmPassword" 
            type="password" 
            required
            value={confirmPassword}
            onChange={(e) => handleConfirmPasswordChange(e.target.value)}
          />
        </div>
        <Button type="submit" className="w-full">
          Sign Up
        </Button>
      </div>
    </form>
  )
}

export {
  SignUpForm,
}