import { Link } from "react-router-dom";
import { Text } from "../styles/primitives/Text";
import { Card } from "../styles/component/Card";
import { Center } from "../styles/layout/Center";
import { Stack } from "../styles/primitives/Stack";

const NotFoundPage = () => {
  return (
    <Center>
      <Card>
        <Stack gap={16}>
          <Text variant="h1">404 — Page not found</Text>

          <Text>
            The page you are looking for does not exist or has been moved.
          </Text>

          <Link to="/">
            Go back to the Home page
          </Link>
        </Stack>
      </Card>
    </Center>
  );
};

export default NotFoundPage;
