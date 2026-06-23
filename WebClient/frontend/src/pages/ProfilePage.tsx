import { Text } from "../styles/primitives/Text";
import { Card } from "../styles/component/Card";
import { Stack } from "../styles/primitives/Stack";
import { useAuth } from "../hooks/useAuth";

const ProfilePage = () => {
  const { state } = useAuth();

  return (
    <Card>
      <Stack gap={16}>
        <Text variant="h1">Profile</Text>
        <Text variant="body">
          Status: {state}
        </Text>
        <Text variant="muted">
          Profile management will be available soon.
        </Text>
      </Stack>
    </Card>
  );
};

export default ProfilePage;
