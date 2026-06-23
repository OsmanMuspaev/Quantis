import { useParams } from "react-router-dom";
import { Text } from "../styles/primitives/Text";
import { Card } from "../styles/component/Card";
import { Stack } from "../styles/primitives/Stack";

const TestResultPage = () => {
  const { id } = useParams<{ id: string }>();

  return (
    <Card>
      <Stack gap={16}>
        <Text variant="h1">Test Results</Text>
        <Text variant="body">
          Test ID: {id}
        </Text>
        <Text variant="muted">
          Detailed results will be available soon.
        </Text>
      </Stack>
    </Card>
  );
};

export default TestResultPage;
