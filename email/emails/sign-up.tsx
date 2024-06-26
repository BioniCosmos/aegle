import {
  Body,
  Button,
  Container,
  Heading,
  Hr,
  Html,
  Link,
  Preview,
  Tailwind,
  Text,
} from '@react-email/components'
import React from 'react'

export default function SignUp() {
  const link = '{{.}}'
  return (
    <Html>
      <Preview>Welcome to Aegle!</Preview>
      <Tailwind>
        <Body>
          <Container className="px-4">
            <Heading>Sign Up Verification</Heading>
            <Button
              className="bg-black rounded text-white font-semibold px-5 py-3"
              href={link}
            >
              Verify
            </Button>
            <Text>or copy and paste this URL into your browser:</Text>
            <Link href={link}>{link}</Link>
            <Text>
              If you didn’t try to sign up, you can safely ignore this email.
            </Text>
            <Hr />
            <Text className="text-gray-400">Aegle</Text>
          </Container>
        </Body>
      </Tailwind>
    </Html>
  )
}
