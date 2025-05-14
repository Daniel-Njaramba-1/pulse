<script lang="ts">
    import { goto } from "$app/navigation";
    import { Button } from "$lib/components/ui/button/index";
    import * as Card from "$lib/components/ui/card/index";
    import * as Alert from "$lib/components/ui/alert/index";
    import { Input } from "$lib/components/ui/input/index";
    import { Label } from "$lib/components/ui/label/index";
    import { authHelpers } from "$lib/stores/auth";
    
    let email = $state<string>("");
    let username = $state<string>("");
    let password = $state<string>("");
    let confirmPassword = $state<string>("");
    let registerLoading = $state<boolean>(false);
    let registerError = $state<string>("");

    async function handleRegister(event:SubmitEvent): Promise<void> {
        event.preventDefault();
        registerLoading = true;
        registerError = '';

        if (!email.includes('@')) {
            registerError = 'Please enter a valid email address';
            registerLoading = false;
            return;
        }

        if (password.length < 6) {
            registerError = 'Password must be at least 6 characters';
            registerLoading = false;
            return;
        }

        if (password !== confirmPassword) {
            registerError = 'Passwords do not match';
            registerLoading = false;
            return;
        }

        try {
            const result:any = await authHelpers.register(email, username, password);
            if (result.success) {
                goto("/");
            } else {
                registerError = result.error;
            }
        } catch (error) {
            registerError = "Error occurred while logging in";
        } finally {
            registerLoading= false;
        }
    }
</script>

<div class="flex justify-center items-center min-h-screen bg-gray-100">
    <Card.Root class="w-full max-w-md shadow-lg">
        <Card.Header class="space-y-1">
            <Card.Title class="text-2xl text-center">Admin Register</Card.Title>
            <Card.Description class="text-center text-gray-500">
                Input details to create your account
            </Card.Description>
        </Card.Header>
        <Card.Content>
            <form onsubmit={handleRegister} class="space-y-4">
                <div class="space-y-2">
                    <Label for="email" class="text-sm font-medium">Email</Label>
                    <Input 
                        id="email"
                        type="email" 
                        bind:value={email} 
                        placeholder=""
                        required
                        class="w-full"
                    />
                </div>
                <div class="space-y-2">
                    <Label for="username" class="text-sm font-medium">Username</Label>
                    <Input 
                        id="username"
                        type="text" 
                        bind:value={username} 
                        placeholder=""
                        required
                        class="w-full"
                    />
                </div>
                <div class="space-y-2">
                    <Label for="password" class="text-sm font-medium">Password</Label>
                    <Input 
                        id="password"
                        type="password" 
                        bind:value={password} 
                        placeholder=""
                        required
                        class="w-full"
                    />
                </div>
                <div class="space-y-2">
                    <Label for="confirmPassword" class="text-sm font-medium">Confirm Password</Label>
                    <Input 
                        id="confirmPassword"
                        type="password" 
                        bind:value={confirmPassword} 
                        placeholder=""
                        required
                        class="w-full"
                    />
                </div>
                {#if registerError}
                    <Alert.Root variant="destructive" class="my-4">
                        <Alert.Title>Authentication Error</Alert.Title> 
                        <Alert.Description>{registerError}</Alert.Description>
                    </Alert.Root>
                {/if}
                <Button 
                    type="submit" 
                    class="w-full" 
                    variant="default"
                    disabled={registerLoading}
                >
                    {registerLoading ? 'Registering...' : 'Register'}
                </Button>
            </form>
        </Card.Content>
        <Card.Footer class="flex justify-center">
            <p class="text-sm text-gray-500">
                Already have an account? <a href="/login" class="text-blue-600 hover:underline">Login</a>
            </p>
        </Card.Footer>
    </Card.Root>
</div>