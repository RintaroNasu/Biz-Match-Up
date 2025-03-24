import { NextRequest, NextResponse } from 'next/server';
import jwt from 'jsonwebtoken';
import prisma from '@/lib/prisma';

const JWT_SECRET = process.env.JWT_SECRET;

export async function POST(req: NextRequest) {
  const { email, password } = await req.json();
  console.log('email:', email);
  console.log('password:', password);
  try {
    const user = await prisma.user.create({
      data: {
        email,
        password,
      },
    });
    console.log('user:', user);
    if (JWT_SECRET) {
      const token = jwt.sign({ id: user.id, email: user.email }, JWT_SECRET, {
        expiresIn: '1h',
      });
      return NextResponse.json(
        { message: 'サインアップ成功', user, token },
        { status: 200 },
      );
    }
    return NextResponse.json(user);
  } catch (e) {
    return NextResponse.json(e);
  }
}
