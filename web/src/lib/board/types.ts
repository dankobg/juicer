export type Color = 'w' | 'b';

export type Row = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7;
export type Col = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7;

export type Rank = 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8;
export type File = 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h';

export type Coordinate = `${File}${Rank}`;

export type PieceSymbol = 'p' | 'n' | 'b' | 'r' | 'q' | 'k';
export type WhitePieceSymbol = Uppercase<PieceSymbol>;
export type BlackPieceSymbol = Lowercase<PieceSymbol>;

export type PromotionPieceSymbol = Exclude<PieceSymbol, 'p' | 'k'>;
export type PromotionWhitePieceSymbol = Uppercase<PromotionPieceSymbol>;
export type PromotionBlackPieceSymbol = Lowercase<PromotionPieceSymbol>;

export enum CastleRights {
	None = 0,
	WhiteKingSide = 1,
	WhiteQueenSide = 2,
	BlackKingSide = 4,
	BlackQueenSide = 8,
}

export type DragPosition = {
	initialX: number;
	initialY: number;
	dx: number;
	dy: number;
};

export type DropOffBoardAction = 'trash' | 'snapback';
