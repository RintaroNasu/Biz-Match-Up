import { NextRequest, NextResponse } from 'next/server';
import prisma from '@/lib/prisma';

export async function POST(
  req: NextRequest,
  context: { params: { userId: string } },
) {
  const body = await req.json();

  const userId = parseInt(context.params.userId);

  try {
    const updatedUser = await prisma.user.update({
      where: { id: userId },
      data: body,
    });

    return NextResponse.json({ success: true, user: updatedUser });
  } catch (error) {
    return NextResponse.json({ success: false, error });
  }
}
