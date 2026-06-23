import { useParams } from "react-router-dom";
import { Text } from "../styles/primitives/Text";
import { Card } from "../styles/component/Card";
import { Stack } from "../styles/primitives/Stack";

const MyTestsPage = () => {
  const { id } = useParams<{ id: string }>();

  return (
    <Card>
      <Stack gap={16}>
        <Text variant="h1">Tests</Text>
        <Text variant="body">
          Course ID: {id}
        </Text>
        <Text variant="muted">
          Test management will be available soon.
        </Text>
      </Stack>
    </Card>
  );
};

export default MyTestsPage;
