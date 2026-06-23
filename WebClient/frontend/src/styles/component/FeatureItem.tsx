import { Card } from './Card';
import { Text } from '../primitives/Text';
import { Stack } from '../primitives/Stack';

interface FeatureItemProps {
  title: string;
  description: string;
}

export const FeatureItem = ({ title, description }: FeatureItemProps) => (
  <Card>
    <Stack gap={8}>
      <Text variant="h3">{title}</Text>
      <Text variant="muted">{description}</Text>
    </Stack>
  </Card>
);
