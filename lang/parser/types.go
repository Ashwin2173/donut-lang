package parser

// primary types
const (
	TypeUnknown TypeId = iota
	TypeVoid
	TypeInt
	TypeFloat
	TypeString
	TypeBoolean
)

// expression operation
const (
	Uop_ineg ExpOpTypeId = iota
	Uop_fneg
	Uop_not
	Bop_none
	Bop_iadd
	Bop_isub
	Bop_imul
	Bop_idiv
	Bop_irem
	Bop_fadd
	Bop_fsub
	Bop_fmul
	Bop_fdiv
	Bop_frem
	Bop_i2f
	Bop_icmp_eq
	Bop_icmp_ne
	Bop_icmp_lt
	Bop_icmp_le
	Bop_icmp_gt
	Bop_icmp_ge
	Bop_fcmp_eq
	Bop_fcmp_ne
	Bop_fcmp_lt
	Bop_fcmp_le
	Bop_fcmp_gt
	Bop_fcmp_ge
	Bop_iand
	Bop_ior
	Bop_s_concat
	Bop_s_eq
	Bop_s_ne
)

type UnaryTypeKey struct {
	operation TokenType
	right     TypeId
}

type BinaryTypeKey struct {
	left      TypeId
	operation TokenType
	right     TypeId
}

type BinaryTypeValue struct {
	result        TypeId
	OperationType ExpOpTypeId
	LeftTo        ExpOpTypeId
	RightTo       ExpOpTypeId
}

var unaryRules = map[UnaryTypeKey]ExpOpTypeId{
	{t_minus, TypeInt}:   Uop_ineg,
	{t_minus, TypeFloat}: Uop_fneg,
	{k_not, TypeBoolean}: Uop_not,
}

var binaryRules = map[BinaryTypeKey]BinaryTypeValue{
	{TypeInt, t_plus, TypeInt}:               {TypeInt, Bop_iadd, Bop_none, Bop_none},
	{TypeInt, t_minus, TypeInt}:              {TypeInt, Bop_isub, Bop_none, Bop_none},
	{TypeInt, t_star, TypeInt}:               {TypeInt, Bop_imul, Bop_none, Bop_none},
	{TypeInt, t_slash, TypeInt}:              {TypeInt, Bop_idiv, Bop_none, Bop_none},
	{TypeInt, t_mod, TypeInt}:                {TypeInt, Bop_irem, Bop_none, Bop_none},
	{TypeInt, t_lesser, TypeInt}:             {TypeBoolean, Bop_icmp_lt, Bop_none, Bop_none},
	{TypeInt, t_greater, TypeInt}:            {TypeBoolean, Bop_icmp_gt, Bop_none, Bop_none},
	{TypeInt, t_lesserEquals, TypeInt}:       {TypeBoolean, Bop_icmp_le, Bop_none, Bop_none},
	{TypeInt, t_greaterEquals, TypeInt}:      {TypeBoolean, Bop_icmp_ge, Bop_none, Bop_none},
	{TypeInt, t_equalEqual, TypeInt}:         {TypeBoolean, Bop_icmp_eq, Bop_none, Bop_none},
	{TypeInt, t_bangEqual, TypeInt}:          {TypeBoolean, Bop_icmp_ne, Bop_none, Bop_none},
	{TypeInt, t_plus, TypeFloat}:             {TypeFloat, Bop_fadd, Bop_i2f, Bop_none},
	{TypeInt, t_minus, TypeFloat}:            {TypeFloat, Bop_fsub, Bop_i2f, Bop_none},
	{TypeInt, t_star, TypeFloat}:             {TypeFloat, Bop_fmul, Bop_i2f, Bop_none},
	{TypeInt, t_slash, TypeFloat}:            {TypeFloat, Bop_fdiv, Bop_i2f, Bop_none},
	{TypeInt, t_mod, TypeFloat}:              {TypeFloat, Bop_frem, Bop_i2f, Bop_none},
	{TypeInt, t_lesser, TypeFloat}:           {TypeBoolean, Bop_fcmp_lt, Bop_i2f, Bop_none},
	{TypeInt, t_greater, TypeFloat}:          {TypeBoolean, Bop_fcmp_gt, Bop_i2f, Bop_none},
	{TypeInt, t_lesserEquals, TypeFloat}:     {TypeBoolean, Bop_fcmp_le, Bop_i2f, Bop_none},
	{TypeInt, t_greaterEquals, TypeFloat}:    {TypeBoolean, Bop_fcmp_ge, Bop_i2f, Bop_none},
	{TypeInt, t_equalEqual, TypeFloat}:       {TypeBoolean, Bop_fcmp_eq, Bop_i2f, Bop_none},
	{TypeInt, t_bangEqual, TypeFloat}:        {TypeBoolean, Bop_fcmp_ne, Bop_i2f, Bop_none},
	{TypeFloat, t_plus, TypeFloat}:           {TypeFloat, Bop_fadd, Bop_none, Bop_none},
	{TypeFloat, t_minus, TypeFloat}:          {TypeFloat, Bop_fsub, Bop_none, Bop_none},
	{TypeFloat, t_star, TypeFloat}:           {TypeFloat, Bop_fmul, Bop_none, Bop_none},
	{TypeFloat, t_slash, TypeFloat}:          {TypeFloat, Bop_fdiv, Bop_none, Bop_none},
	{TypeFloat, t_mod, TypeFloat}:            {TypeFloat, Bop_frem, Bop_none, Bop_none},
	{TypeFloat, t_lesser, TypeFloat}:         {TypeBoolean, Bop_fcmp_lt, Bop_none, Bop_none},
	{TypeFloat, t_greater, TypeFloat}:        {TypeBoolean, Bop_fcmp_gt, Bop_none, Bop_none},
	{TypeFloat, t_lesserEquals, TypeFloat}:   {TypeBoolean, Bop_fcmp_le, Bop_none, Bop_none},
	{TypeFloat, t_greaterEquals, TypeFloat}:  {TypeBoolean, Bop_fcmp_ge, Bop_none, Bop_none},
	{TypeFloat, t_equalEqual, TypeFloat}:     {TypeBoolean, Bop_fcmp_eq, Bop_none, Bop_none},
	{TypeFloat, t_bangEqual, TypeFloat}:      {TypeBoolean, Bop_fcmp_ne, Bop_none, Bop_none},
	{TypeFloat, t_plus, TypeInt}:             {TypeFloat, Bop_fadd, Bop_none, Bop_i2f},
	{TypeFloat, t_minus, TypeInt}:            {TypeFloat, Bop_fsub, Bop_none, Bop_i2f},
	{TypeFloat, t_star, TypeInt}:             {TypeFloat, Bop_fmul, Bop_none, Bop_i2f},
	{TypeFloat, t_slash, TypeInt}:            {TypeFloat, Bop_fdiv, Bop_none, Bop_i2f},
	{TypeFloat, t_mod, TypeInt}:              {TypeFloat, Bop_frem, Bop_none, Bop_i2f},
	{TypeFloat, t_lesser, TypeInt}:           {TypeBoolean, Bop_fcmp_lt, Bop_none, Bop_i2f},
	{TypeFloat, t_greater, TypeInt}:          {TypeBoolean, Bop_fcmp_gt, Bop_none, Bop_i2f},
	{TypeFloat, t_lesserEquals, TypeInt}:     {TypeBoolean, Bop_fcmp_le, Bop_none, Bop_i2f},
	{TypeFloat, t_greaterEquals, TypeInt}:    {TypeBoolean, Bop_fcmp_ge, Bop_none, Bop_i2f},
	{TypeFloat, t_equalEqual, TypeInt}:       {TypeBoolean, Bop_fcmp_eq, Bop_none, Bop_i2f},
	{TypeFloat, t_bangEqual, TypeInt}:        {TypeBoolean, Bop_fcmp_ne, Bop_none, Bop_i2f},
	{TypeString, t_plus, TypeString}:         {TypeString, Bop_s_concat, Bop_none, Bop_none},
	{TypeString, t_equalEqual, TypeString}:   {TypeBoolean, Bop_s_eq, Bop_none, Bop_none},
	{TypeString, t_bangEqual, TypeString}:    {TypeBoolean, Bop_s_ne, Bop_none, Bop_none},
	{TypeBoolean, k_and, TypeBoolean}:        {TypeBoolean, Bop_iand, Bop_none, Bop_none},
	{TypeBoolean, k_or, TypeBoolean}:         {TypeBoolean, Bop_ior, Bop_none, Bop_none},
	{TypeBoolean, t_equalEqual, TypeBoolean}: {TypeBoolean, Bop_icmp_eq, Bop_none, Bop_none},
	{TypeBoolean, t_bangEqual, TypeBoolean}:  {TypeBoolean, Bop_icmp_ne, Bop_none, Bop_none},
}
