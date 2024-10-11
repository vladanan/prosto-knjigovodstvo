// daje random int uključujući min i max vrednosti
function rndInt(min, max) {
	return Math.floor(Math.random() * (max - min + 1) ) + min;
}

// daje novi niz sa nasumičnim elementima iz celog niza
// menja niz koji je ubačen!!!
const array_rand2 = (niz) => {
	let novi_niz = [0];
	let i = 0;
	let izabrani_clan;
	while (niz.length > 0) {
		izabrani_clan = rndInt(0, niz.length-1);
		novi_niz[i] = niz[izabrani_clan];
		niz.splice(izabrani_clan, 1);
		i++;
	}
	return novi_niz;
}

// daje nov niz sa nasumičnim elementima iz niza tako da ih ima koliko broj_clanova
// PRAVI DUPLE I MORA DA SE SREDI PO UGLEDU NA splice MEHANIZAM ONOG GORE!!!!
function array_rand(niz, broj_clanova) {
	//console.log('niz:' + niz, 'broj:' + broj_clanova);
	//console.log(niz.length);
	let novi_niz = [0];
	for (let i = 0; i < broj_clanova; i++) {
		let clan = rndInt(0, niz.length-1);
		//console.log('clan: ' + clan);
		//console.log(niz[clan]);
		novi_niz[i] = niz[clan];
		//console.log('novi niz for:'+novi_niz);
	}
	//console.log('novi niz out:'+novi_niz);
	return novi_niz;
}

//lažni shuffle samo za nizove sa 10 članova
function shuffle() {
	let kombinacija = rndInt(0, 2);
	let novi_niz = [0];
	let kombinacije = [
		[4, 8, 7, 0, 2, 5, 3, 1, 9, 6],
		[6, 9, 1, 3, 5, 2, 0, 7, 8, 4],
		[2, 9, 3, 8, 6, 5, 1, 0, 4, 7]
	];
	let izabrana = kombinacije[kombinacija];

	//for (let i = 0; i < 10; i++) {
		novi_niz = izabrana;
	//}
	return novi_niz;
}


function O2m_1_100 (props) {

	let pitanje = 'neko pitanje';

	//console.log('props', props)

	let a = props.novo_pitanje.a;
	let op2 = props.novo_pitanje.op2;
	let b = props.novo_pitanje.b;

	pitanje =
	<>

		<label className='text-2xl ml-2' htmlFor="r">{a} {op2} {b} = </label>
		<input onChange={(e) => props.setOdgovor(e.target.value)} className='text-2xl w-16  bg-lime-300 px-1 border border-lime-700 rounded-md' type='number' min='0' max='100' id='r' name='r' defaultValue={0} />
		<span className='text-2xl ml-2 font-bold text-green-700 '>&nbsp;&nbsp;{props.resenje}</span>
	</>

	let o2m_1_100 = <>

		<div
		className="p-1 border border-black rounded-md bg-gradient-to-t from-lime-100 via-white to-white shadow-xl">
				{pitanje}
		</div>

		<div hidden={!props.vidljiv_odgovor}
		className="p-1 h-fit border border-t-0 border-black rounded-md bg-lime-50 shadow-xl">
				{props.rezultat}
		</div>

	</>

	return (
		o2m_1_100
	)
}


function O2m_1_100txt (props) {

	let pitanje = 'neko pitanje';

	//console.log('props', props)

	let a = props.novo_pitanje.a;
	//let op2 = props.novo_pitanje.op2;
	let b = props.novo_pitanje.b;
	let t11 = props.novo_pitanje.t11;
	let t22 = props.novo_pitanje.t22;
	let t33 = props.novo_pitanje.t33;

	let c;
	let label;

	if (a > b) {
		c = a - b;
		label = t11 + props.i18n.Rec.Ima + a + " " + t33 + props.i18n.Rec.Aa + t22 + " " + c + props.i18n.Rec.Manje_kraj;
	} else {
		c = b - a;
		label = t11 + props.i18n.Rec.Ima + a + " " + t33 + props.i18n.Rec.Aa + t22 + " " + c + props.i18n.Rec.Vise_kraj;
	}

	pitanje =
	<>

		<label className='text-2xl' htmlFor="r">{label} </label>
		<input onChange={(e) => props.setOdgovor(e.target.value)} className='text-2xl w-16 bg-lime-300 px-1 border border-lime-700 rounded-md' type='number' min='0' max='100' id='r' name='r' defaultValue={0} />
		<span className='text-2xl ml-2 font-bold text-green-700 '>&nbsp;&nbsp;{props.resenje}</span>
	</>

	let o2m_1_100txt = <>

		<div
		className="p-1 border border-black rounded-md bg-gradient-to-t from-lime-100 via-white to-white shadow-xl">
				{pitanje}
		</div>

		<div hidden={!props.vidljiv_odgovor}
		className="p-1 h-fit border border-t-0 border-black rounded-md bg-lime-50 shadow-xl">
				{props.rezultat}
		</div>

	</>

	return (
		o2m_1_100txt
	)
}



function O2m_mnozenje (props) {

	let pitanje = 'neko pitanje';

	//console.log('mnozenje props', props)

	let a = props.novo_pitanje.a;
	let op2 = props.novo_pitanje.op2;
	let b = props.novo_pitanje.b;

	//console.log('mnozenje vars', a,op2,b)

	pitanje =
	<>
		<label className='text-2xl ml-2' htmlFor="r">{a} {op2} {b} = </label>
		<input onChange={(e) => props.setOdgovor(e.target.value)} className='text-2xl w-20 bg-lime-300 px-1 border border-lime-700 rounded-md' type='number' min='1' max='100' id='r' name='r' defaultValue={0} />
		<span className='text-2xl ml-2 font-bold text-green-700 '>&nbsp;&nbsp;{props.resenje}</span>
	</>

	let o2m_mnozenje = <>

		<div
		className="p-1 border border-black rounded-md bg-gradient-to-t from-lime-100 via-white to-white shadow-xl">
				{pitanje}
		</div>

		<div hidden={!props.vidljiv_odgovor}
		className="p-1 h-fit border border-t-0 border-black rounded-md bg-lime-50 shadow-xl">
				{props.rezultat}
		</div>

	</>

	return (
		o2m_mnozenje
	)
}


export function reactRenderZadaci_o2(translations) {
	// console.log("json translations:", JSON.parse(translations))
	const rootElement = document.getElementById("root");
	if (!rootElement) {
		throw new Error(`Could not find element with iidd ${"root"}`);
	}
	const reactRoot = ReactDOM.createRoot(rootElement);
	reactRoot.render(<Zadaci_o2 i18n={JSON.parse(translations)}/>);
}


function Zadaci_o2 (props) {

	//console.log('PROPOVI: ', props)

	const [zadatak, setZadatak] = React.useState('o2m_1_100');
	const [vidljiv_odgovor, setVidljiv_odgovor] = React.useState(false);
	const [novo_pitanje, setNovo_pitanje] = React.useState({
		a: 2,
		op2: '+',
		b: 2,
		t11: props.i18n.Novo_pitanje.T11,
		t22: props.i18n.Novo_pitanje.T22,
		t33: props.i18n.Novo_pitanje.T33,
		kombi: ''
	});

	const [odgovor, setOdgovor] = React.useState(null);
	const [resenje, setResenje] = React.useState(null);
	const [rezultat, setRezultat] = React.useState(null);

	const [prviCinilac, setPrviCinilac] = React.useState({
		psvi: true,
		p1: true,
		p2: true,
		p3: true,
		p4: true,
		p5: true,
		p6: true,
		p7: true,
		p8: true,
		p9: true
	})

	const [drugiCinilac, setDrugiCinilac] = React.useState({
		dsvi: true,
		d1: true,
		d2: true,
		d3: true,
		d4: true,
		d5: true,
		d6: true,
		d7: true,
		d8: true,
		d9: true
	})

	const promeniZadatak = (zadatak_pt) => {
		setZadatak(zadatak_pt);

		if (zadatak_pt == 'o2m_1_100' || zadatak_pt == 'o2m_1_100txt') {
			setNovo_pitanje({
				a: 0,
				op2: '+',
				b: 0,
				t11: props.i18n.SetNovo_pitanje.T11,
				t22: props.i18n.SetNovo_pitanje.T22,
				t33: props.i18n.SetNovo_pitanje.T33,
				kombi: ''
			});
		} else {
			setNovo_pitanje({
				a: 0,
				op2: '*',
				b: 0,
				t11: props.i18n.SetNovo_pitanje.T11,
				t22: props.i18n.SetNovo_pitanje.T22,
				t33: props.i18n.SetNovo_pitanje.T33,
				kombi: ''
			});
		}

	}

	const vidiOdgovor = () => {

		let rezultat;

		let yes_niz = [
			'/y1.webp',
			'/y2.gif',
			'/y3.png',
			'/y4.webp',
			'/y5.gif',
			'/y6.gif',
			'/y7.gif',
			'/y8.jpeg',
			'/y9.jpg',
			'/y10.gif'
		]

		let no_niz = [
			'/n1.webp',
			'/n2.gif',
			'/n3.webp',
			'/n4.jpeg',
			'/n5.png',
			'/n6.jpg',
			'/n7.gif',
			'/n8.gif',
			'/n9.gif',
			'/n10.jpg'
		]

		let rnd_yes_niz = rndInt(0, 10)
		let rnd_no_niz = rndInt(0, 10)

		switch (zadatak) {
			case 'o2m_1_100':
				var {a, op2, b} = novo_pitanje;
				var r = odgovor;

				//console.log(a, op2, b, r)

				var rr = "";

				var tacno = 0;
				if (op2 == '-') {
					var rr = a - b;
				} else {
					var rr = a + b;
				}

				if (rr == r) {
				tacno = 1;
				} else {
				tacno = 0;
				}

				setResenje(rr);

				if (tacno == 1) {
					rezultat = <>
						<p style={{textAlign: "center", background: "skyblue", fontSize: "30px"}}>
							&#10004;
						</p>
						{/*
						<img src={yes_niz[rnd_yes_niz]} />
						<audio controls autoPlay>
							<source src='/slavuj2.mp3' type='audio/mpeg' />
						</audio> */}
					</>
				} else {
					rezultat = <>
						<p style={{background: "red", textAlign: "center", fontSize: "30px"}}>
							&#10008;
						</p>
						{/*
						<img src={no_niz[rnd_no_niz]} />
						<audio controls autoPlay>
							<source src='/beba.mp3' type='audio/mpeg' />
						</audio> */}
					</>
				}
				break;

			case 'o2m_1_100txt':
				var {a, b} = novo_pitanje;
				var r = odgovor;

				var rr = a + b;
				var tacno = 0;

				if (rr == r) {
					tacno = 1;
				} else {
					tacno = 0;
				}

				setResenje(rr);

				if (tacno == 1) {
					rezultat = <>
						<p style={{textAlign: "center", background: "skyblue", fontSize: "30px"}}>
							&#10004;
						</p>
						{/*
						<img src={yes_niz[rnd_yes_niz]} />
						<audio controls autoPlay>
							<source src='/slavuj2.mp3' type='audio/mpeg' />
						</audio> */}
					</>
				} else {
					rezultat = <>
						<p style={{background: "red", textAlign: "center", fontSize: "30px"}}>
							&#10008;
						</p>
						{/*
						<img src={no_niz[rnd_no_niz]} />
						<audio controls autoPlay>
							<source src='/beba.mp3' type='audio/mpeg' />
						</audio> */}
					</>
				}
				break;

			case 'o2m_mnozenje':
				var {a, op2, b} = novo_pitanje;
				var r = odgovor;

				//console.log('odgovor: ', a, op2, b, r)

				var rr = "";

				var tacno = 0;
				if (op2 == '*') {
					var rr = a * b;
				} else {
					var rr = a * b;
				}

				if (rr == r) {
				tacno = 1;
				} else {
				tacno = 0;
				}

				setResenje(rr);

				if (tacno == 1) {
					rezultat = <>
						<p style={{textAlign: "center", background: "skyblue", fontSize: "30px"}}>
							&#10004;
						</p>
						{/*
						<img src={yes_niz[rnd_yes_niz]} />
						<audio controls autoPlay>
							<source src='/slavuj2.mp3' type='audio/mpeg' />
						</audio> */}
					</>
				} else {
					rezultat = <>
						<p style={{background: "red", textAlign: "center", fontSize: "30px"}}>
							&#10008;
						</p>
						{/*
						<img src={no_niz[rnd_no_niz]} />
						<audio controls autoPlay>
							<source src='/beba.mp3' type='audio/mpeg' />
						</audio> */}
					</>
				}
				break;
			default:
		}



		if (vidljiv_odgovor == true) {
			setVidljiv_odgovor(false);
		} else {
			setVidljiv_odgovor(true);
			setRezultat(rezultat);
		}
		//console.log('odgovor: ', odgovor)
	}

	const novoPitanje = () => {

		//console.log(zadatak)

		switch (zadatak) {
			case 'o2m_1_100':
				let a = rndInt(2, 99);
				let op1 = rndInt(0, 1);
				let op2 = '';
				let b = rndInt(2, 99);

				if (op1 == 0) {
					if (a < b) {
						let tmp = a;
						a = b;
						b = tmp;
					}
				} else {
					while (a+b > 100) {
						a = rndInt(2, 99);
						b = rndInt(2, 99);
					}
				}

				if (op1 == 0) {
					op2 = "-";
				} else {
					op2 = "+";
				}

				setResenje(null);

				setNovo_pitanje({
					a,
					op2,
					b,
					t11: props.i18n.SetNovo_pitanje.T11,
					t22: props.i18n.SetNovo_pitanje.T22,
					t33: props.i18n.SetNovo_pitanje.T33,
					kombi: ''
				});
				break;

			case 'o2m_1_100txt':
				let aa = rndInt(2, 99);
				let bb = rndInt(2, 99);

				let t1 = [
					props.i18n.Novo_pitanjeT1[0],
					props.i18n.Novo_pitanjeT1[1],
					props.i18n.Novo_pitanjeT1[2],
					props.i18n.Novo_pitanjeT1[3],
					props.i18n.Novo_pitanjeT1[4],
					props.i18n.Novo_pitanjeT1[5],
					props.i18n.Novo_pitanjeT1[6],
					props.i18n.Novo_pitanjeT1[7],
					props.i18n.Novo_pitanjeT1[8],
					props.i18n.Novo_pitanjeT1[9],
				];

				let t2 = [
					props.i18n.Novo_pitanjeT2[0],
					props.i18n.Novo_pitanjeT2[1],
					props.i18n.Novo_pitanjeT2[2],
					props.i18n.Novo_pitanjeT2[3],
					props.i18n.Novo_pitanjeT2[4],
					props.i18n.Novo_pitanjeT2[5],
					props.i18n.Novo_pitanjeT2[6],
					props.i18n.Novo_pitanjeT2[7],
					props.i18n.Novo_pitanjeT2[8],
					props.i18n.Novo_pitanjeT2[9],
				];

				let t3 = [
					props.i18n.Novo_pitanjeT3[0],
					props.i18n.Novo_pitanjeT3[1],
					props.i18n.Novo_pitanjeT3[2],
					props.i18n.Novo_pitanjeT3[3],
					props.i18n.Novo_pitanjeT3[4],
					props.i18n.Novo_pitanjeT3[5],
					props.i18n.Novo_pitanjeT3[6],
					props.i18n.Novo_pitanjeT3[7],
					props.i18n.Novo_pitanjeT3[8],
					props.i18n.Novo_pitanjeT3[9],
				];

				// let t1 = ["Марко", "Mајмун", "Милош", "Kрокодил", "Стефан", "Слепи миш", "Илија", "Слон", "Јован", "Горила"];

				// let t2 = ["Маша", "жирафа", "Јелена", "хијена", "Ивана", "чапља", "Милена", "кокошка", "Наташа", "овца"];

				// let t3 = ["камиончића", "ексера", "колача", "чекића", "динара", "фломастера" , "лизала", "цветића", "лептирића", "ћевапчића"];

				let t11 = t1[rndInt(0, 9)];
				let t22 = t2[rndInt(0, 9)];
				let t33 = t3[rndInt(0, 9)];

				while (aa+bb > 100) {
					aa = rndInt(2, 99);
					bb = rndInt(2, 99);
				}

				setResenje(null);

				setNovo_pitanje({
					a: aa,
					op2: '+', //ovde je nepotrebno ali da ne pravi neke bagove
					b: bb,
					t11,
					t22,
					t33,
					kombi: ''
				});
				break;

			case 'o2m_mnozenje':

				//console.log(prviCinilac)
				// punjenje niza (potrebnog za poređenje za random) na osnovu stanja prviCinilas
				let prvi_niz = []
				let i = 0
				let falses = 0
				for (let cinilac in prviCinilac) {
					if (prviCinilac[cinilac] == true) {
						prvi_niz.push(i)
					} else {
						falses++
					}
					i++
				}
				//console.log(prvi_niz)

				if (falses == 10) { // zaštita od toga da nijedno polje za prvi činilac nije odabrano čime while ide u infinite loop
					prviCinilac.p1 = true
					prvi_niz[1] = 1
				}

				let am

				while (!prvi_niz.includes(am)) {
					am = rndInt(1, 9);
				}

				//let op1m = Math.floor(Math.random() * 2);
				let op2m = '*';

				let drugi_niz = []
				let i2 = 0
				let falses2 = 0
				for (let cinilac in drugiCinilac) {
					if (drugiCinilac[cinilac] == true) {
						drugi_niz.push(i2)
					} else {
						falses2++
					}
					i2++
				}
				//console.log(drugi_niz)

				if (falses2 == 10) { // zaštita od toga da nijedno polje za prvi činilac nije odabrano čime while ide u infinite loop
					drugiCinilac.d1 = true
					drugi_niz[1] = 1
				}

				let bm

				while (!drugi_niz.includes(bm)) {
					bm = rndInt(1, 9);
				}


				//let bm = Math.floor(Math.random() * 11);
				am == 0 ? am = 1 : am = am
				bm == 0 ? bm = 1 : bm = bm

				//console.log(am, op2m, bm)

				setResenje(null);

				setNovo_pitanje({
					a: am,
					op2: op2m,
					b: bm,
					t11: props.i18n.SetNovo_pitanje.T11,
					t22: props.i18n.SetNovo_pitanje.T22,
					t33: props.i18n.SetNovo_pitanje.T33,
					kombi: ''
				});
				break;

			default:
				break;
		}

		setVidljiv_odgovor(false);
		//setOdgovor(null);

	}


	return (
		<>
			<p className="text-2xl mt-5 ml-2">{props.i18n.Title}</p>

			<div className="mx-2 mt-1 text-black py-3 p-2 max-w-sm border-2 border-sky-600 rounded-md bg-gradient-to-r from-blue-100 via-lime-100 to-transparent shadow-lg shadow-sky-900">


				<label className={zadatak == 'o2m_1_100'
					? "relative p-1 border-2 border-black rounded-md bg-gradient-to-r from-sky-400 via-lime-400 to-transparent"
					: "relative p-1 border-2 border-black rounded-md bg-gray-400"}
				htmlFor="o2m_1_100">{props.i18n.Zadaci_1_100}
				</label>
				<input
				onClick={(e) => promeniZadatak('o2m_1_100')}
				className="relative ml-2 mt-2"
				type="radio" id="o2m_1_100" name="radio_z" value="o2m_1_100"
				defaultChecked
				/>
				<br/>

				<label className={zadatak == 'o2m_1_100txt'
					? "relative p-1 border-2 border-black rounded-md bg-gradient-to-br from-red-500 via-yellow-100 to-green-300"
					: "relative p-1 border-2 border-black rounded-md bg-gray-400"}
				htmlFor="o2m_1_100txt">{props.i18n.Zadaci_smesni_1_100}
				</label>
				<input
				onClick={(e) => promeniZadatak('o2m_1_100txt')}
				className="relative ml-2 mt-5"
				type="radio" id="o2m_1_100txt" name="radio_z" value="o2m_1_100txt"
				/>
				<br/>

				<label className={zadatak == 'o2m_mnozenje'
					? "relative p-1 border-2 border-black rounded-md bg-gradient-to-r from-sky-400 via-lime-400 to-transparent"
					: "relative p-1 border-2 border-black rounded-md bg-gray-400"} htmlFor="o2m_mnozenje">{props.i18n.Mnozenje_1_100}
				</label>
				<input
				onClick={(e) => promeniZadatak('o2m_mnozenje')}
				className="relative ml-2 mt-5"
				type="radio" id="o2m_mnozenje" name="radio_z" value="o2m_mnozenje"
				/>
				<br/>


				{zadatak == 'o2m_1_100' && (
				<>
					<br/>
					<button onClick={vidiOdgovor}
					className="ml-10 my-2 relative
					border-2 border-gray-500 rounded-md p-1 bg-lime-50 hover:bg-lime-200" >
						{props.i18n.Result}</button>

					<button onClick={novoPitanje}
					className="ml-8 my-2 relative
					border-2 border-gray-500 rounded-md p-1 bg-blue-300 hover:bg-blue-400">
						{props.i18n.New_a}</button>
					<br/>

					<O2m_1_100
						novo_pitanje={novo_pitanje}
						vidljiv_odgovor={vidljiv_odgovor}
						setOdgovor={setOdgovor}
						rezultat={rezultat}
						resenje={resenje}
					/>
				</>
				)}

				{zadatak == 'o2m_1_100txt' && (
				<>

					<br/>
					<button onClick={vidiOdgovor}
					className="ml-10 my-2 relative
					border-2 border-gray-500 rounded-md p-1 bg-lime-50 hover:bg-lime-200" >
						{props.i18n.Result}</button>

					<button onClick={novoPitanje}
					className="ml-8 my-2 relative
					border-2 border-gray-500 rounded-md p-1 bg-blue-300 hover:bg-blue-400">
						{props.i18n.New_a}</button>
					<br/>

					<O2m_1_100txt
						novo_pitanje={novo_pitanje}
						vidljiv_odgovor={vidljiv_odgovor}
						setOdgovor={setOdgovor}
						rezultat={rezultat}
						resenje={resenje}
						i18n={props.i18n}
					/>
				</>
				)}


				{zadatak == 'o2m_mnozenje' && (
				<>

					<div className='mt-6 grid grid-rows-6 gap-1'>

						<p className="p-1 border-2 border-black rounded-md bg-lime-100 hover:bg-violet-100"
						>{props.i18n.Mnozenje_1_100_1}</p>

						<div className='grid grid-cols-10 gap-0'>

							<label className="w-6 text-center border-x border-t-2 border-black rounded-md bg-lime-300 hover:bg-violet-200"
							htmlFor="psvi">{props.i18n.Mnozenje_1_100_svi}</label>

							<label className="w-6 text-center border-x border-t-2 border-black rounded-md bg-lime-100 hover:bg-violet-200"
							htmlFor="p1">1</label>

							<label className="w-6 text-center border-x border-t-2 border-black rounded-md bg-lime-100 hover:bg-violet-200"
							htmlFor="p2">2</label>

							<label className="w-6 text-center border-x border-t-2 border-black rounded-md bg-lime-100 hover:bg-violet-200"
							htmlFor="p3">3</label>

							<label className="w-6 text-center border-x border-t-2 border-black rounded-md bg-lime-100 hover:bg-violet-200"
							htmlFor="p4">4</label>

							<label className="w-6 text-center border-x border-t-2 border-black rounded-md bg-lime-100 hover:bg-violet-200"
							htmlFor="p5">5</label>

							<label className="w-6 text-center border-x border-t-2 border-black rounded-md bg-lime-100 hover:bg-violet-200"
							htmlFor="p6">6</label>

							<label className="w-6 text-center border-x border-t-2 border-black rounded-md bg-lime-100 hover:bg-violet-200"
							htmlFor="p7">7</label>

							<label className="w-6 text-center border-x border-t-2 border-black rounded-md bg-lime-100 hover:bg-violet-200"
							htmlFor="p8">8</label>

							<label className="w-6 text-center border-x border-t-2 border-black rounded-md bg-lime-100 hover:bg-violet-200"
							htmlFor="p9">9</label>

						</div>

						<div className='ml-1 grid grid-cols-10 gap-1'>

							<input
							onChange={(e) => {
								//console.log('before change checked: ', e.target.checked)
								if (prviCinilac.psvi == true) {
									setPrviCinilac({
										psvi: false,
										p1: true,
										p2: false,
										p3: false,
										p4: false,
										p5: false,
										p6: false,
										p7: false,
										p8: false,
										p9: false
									})
								} else {
									setPrviCinilac({
										psvi: true,
										p1: true,
										p2: true,
										p3: true,
										p4: true,
										p5: true,
										p6: true,
										p7: true,
										p8: true,
										p9: true
									})
								}
								//console.log('after change checked?: ', e.target.checked)
							}}
							className="w-fit"
							type="checkbox"
							id="psvi"
							name="p"
							value="svi"
							checked={prviCinilac.psvi == true ? true : false}
							/>

							<input
							onChange={(e) => {
								if (prviCinilac.p1 == true) {
									setPrviCinilac(prviCinilac => ({
										...prviCinilac,
										...{p1: false}
									}))
								} else {
									setPrviCinilac(prviCinilac => ({
										...prviCinilac,
										...{p1: true}
									}))
								}
							}}
							className="w-fit"
							type="checkbox"
							id="p1"
							name="p"
							value="1"
							checked={prviCinilac.p1 == true ? true : false}
							/>

							<input
							onChange={(e) => {
								if (prviCinilac.p2 == true) {
									setPrviCinilac(prviCinilac => ({
										...prviCinilac,
										...{p2: false}
									}))
								} else {
									setPrviCinilac(prviCinilac => ({
										...prviCinilac,
										...{p2: true}
									}))
								}
							}}
							className="w-fit"
							type="checkbox"
							id="p2"
							name="p"
							value="2"
							checked={prviCinilac.p2 == true ? true : false}
							/>

							<input
							onChange={(e) => {
								if (prviCinilac.p3 == true) {
									setPrviCinilac(prviCinilac => ({
										...prviCinilac,
										...{p3: false}
									}))
								} else {
									setPrviCinilac(prviCinilac => ({
										...prviCinilac,
										...{p3: true}
									}))
								}
							}}
							className="w-fit"
							type="checkbox"
							id="p3"
							name="p"
							value="3"
							checked={prviCinilac.p3 == true ? true : false}
							/>

							<input
							onChange={(e) => {
								if (prviCinilac.p4 == true) {
									setPrviCinilac(prviCinilac => ({
										...prviCinilac,
										...{p4: false}
									}))
								} else {
									setPrviCinilac(prviCinilac => ({
										...prviCinilac,
										...{p4: true}
									}))
								}
							}}
							className="w-fit"
							type="checkbox"
							id="p4"
							name="p"
							value="4"
							checked={prviCinilac.p4 == true ? true : false}
							/>

							<input
							onChange={(e) => {
								if (prviCinilac.p5 == true) {
									setPrviCinilac(prviCinilac => ({
										...prviCinilac,
										...{p5: false}
									}))
								} else {
									setPrviCinilac(prviCinilac => ({
										...prviCinilac,
										...{p5: true}
									}))
								}
							}}
							className="w-fit"
							type="checkbox"
							id="p5"
							name="p"
							value="5"
							checked={prviCinilac.p5 == true ? true : false}
							/>

							<input
							onChange={(e) => {
								if (prviCinilac.p6 == true) {
									setPrviCinilac(prviCinilac => ({
										...prviCinilac,
										...{p6: false}
									}))
								} else {
									setPrviCinilac(prviCinilac => ({
										...prviCinilac,
										...{p6: true}
									}))
								}
							}}
							className="w-fit"
							type="checkbox"
							id="p6"
							name="p"
							value="6"
							checked={prviCinilac.p6 == true ? true : false}
							/>

							<input
							onChange={(e) => {
								if (prviCinilac.p7 == true) {
									setPrviCinilac(prviCinilac => ({
										...prviCinilac,
										...{p7: false}
									}))
								} else {
									setPrviCinilac(prviCinilac => ({
										...prviCinilac,
										...{p7: true}
									}))
								}
							}}
							className="w-fit"
							type="checkbox"
							id="p7"
							name="p"
							value="7"
							checked={prviCinilac.p7 == true ? true : false}
							/>

							<input
							onChange={(e) => {
								if (prviCinilac.p8 == true) {
									setPrviCinilac(prviCinilac => ({
										...prviCinilac,
										...{p8: false}
									}))
								} else {
									setPrviCinilac(prviCinilac => ({
										...prviCinilac,
										...{p8: true}
									}))
								}
							}}
							className="w-fit"
							type="checkbox"
							id="p8"
							name="p"
							value="8"
							checked={prviCinilac.p8 == true ? true : false}
							/>

							<input
							onChange={(e) => {
								if (prviCinilac.p9 == true) {
									setPrviCinilac(prviCinilac => ({
										...prviCinilac,
										...{p9: false}
									}))
								} else {
									setPrviCinilac(prviCinilac => ({
										...prviCinilac,
										...{p9: true}
									}))
								}
							}}
							className="w-fit"
							type="checkbox"
							id="p9"
							name="p"
							value="9"
							checked={prviCinilac.p9 == true ? true : false}
							/>

						</div>

						<p className="p-1 border-2 border-black rounded-md bg-sky-100 hover:bg-violet-100"
						>{props.i18n.Mnozenje_1_100_2}</p>

						<div className='grid grid-cols-10 gap-0'>

							<label className="w-6 text-center border-x border-t-2 border-black rounded-md bg-sky-300 hover:bg-violet-200"
							htmlFor="dsvi">{props.i18n.Mnozenje_1_100_svi}</label>

							<label className="w-6 text-center border-x border-t-2 border-black rounded-md bg-sky-100 hover:bg-violet-200"
							htmlFor="d1">1</label>

							<label className="w-6 text-center border-x border-t-2 border-black rounded-md bg-sky-100 hover:bg-violet-200"
							htmlFor="d2">2</label>

							<label className="w-6 text-center border-x border-t-2 border-black rounded-md bg-sky-100 hover:bg-violet-200"
							htmlFor="d3">3</label>

							<label className="w-6 text-center border-x border-t-2 border-black rounded-md bg-sky-100 hover:bg-violet-200"
							htmlFor="d4">4</label>

							<label className="w-6 text-center border-x border-t-2 border-black rounded-md bg-sky-100 hover:bg-violet-200"
							htmlFor="d5">5</label>

							<label className="w-6 text-center border-x border-t-2 border-black rounded-md bg-sky-100 hover:bg-violet-200"
							htmlFor="d6">6</label>

							<label className="w-6 text-center border-x border-t-2 border-black rounded-md bg-sky-100 hover:bg-violet-200"
							htmlFor="d7">7</label>

							<label className="w-6 text-center border-x border-t-2 border-black rounded-md bg-sky-100 hover:bg-violet-200"
							htmlFor="d8">8</label>

							<label className="w-6 text-center border-x border-t-2 border-black rounded-md bg-sky-100 hover:bg-violet-200"
							htmlFor="d9">9</label>

						</div>

						<div className='ml-1 grid grid-cols-10 gap-1'>

							<input
							onChange={(e) => {
								//console.log('before change checked: ', e.target.checked)
								if (drugiCinilac.dsvi == true) {
									setDrugiCinilac({
										dsvi: false,
										d1: true,
										d2: false,
										d3: false,
										d4: false,
										d5: false,
										d6: false,
										d7: false,
										d8: false,
										d9: false
									})
								} else {
									setDrugiCinilac({
										dsvi: true,
										d1: true,
										d2: true,
										d3: true,
										d4: true,
										d5: true,
										d6: true,
										d7: true,
										d8: true,
										d9: true
									})
								}
								//console.log('after change checked?: ', e.target.checked)
							}}
							className="w-fit"
							type="checkbox"
							id="dsvi"
							name="d"
							value="svi"
							checked={drugiCinilac.dsvi == true ? true : false}
							/>

							<input
							onChange={(e) => {
								if (drugiCinilac.d1 == true) {
									setDrugiCinilac(drugiCinilac => ({
										...drugiCinilac,
										...{d1: false}
									}))
								} else {
									setDrugiCinilac(drugiCinilac => ({
										...drugiCinilac,
										...{d1: true}
									}))
								}
							}}
							className="w-fit"
							type="checkbox"
							id="d1"
							name="d"
							value="1"
							checked={drugiCinilac.d1 == true ? true : false}
							/>

							<input
							onChange={(e) => {
								if (drugiCinilac.d2 == true) {
									setDrugiCinilac(drugiCinilac => ({
										...drugiCinilac,
										...{d2: false}
									}))
								} else {
									setDrugiCinilac(drugiCinilac => ({
										...drugiCinilac,
										...{d2: true}
									}))
								}
							}}
							className="w-fit"
							type="checkbox"
							id="d2"
							name="d"
							value="2"
							checked={drugiCinilac.d2 == true ? true : false}
							/>

							<input
							onChange={(e) => {
								if (drugiCinilac.d3 == true) {
									setDrugiCinilac(drugiCinilac => ({
										...drugiCinilac,
										...{d3: false}
									}))
								} else {
									setDrugiCinilac(drugiCinilac => ({
										...drugiCinilac,
										...{d3: true}
									}))
								}
							}}
							className="w-fit"
							type="checkbox"
							id="d3"
							name="d"
							value="3"
							checked={drugiCinilac.d3 == true ? true : false}
							/>

							<input
							onChange={(e) => {
								if (drugiCinilac.d4 == true) {
									setDrugiCinilac(drugiCinilac => ({
										...drugiCinilac,
										...{d4: false}
									}))
								} else {
									setDrugiCinilac(drugiCinilac => ({
										...drugiCinilac,
										...{d4: true}
									}))
								}
							}}
							className="w-fit"
							type="checkbox"
							id="d4"
							name="d"
							value="4"
							checked={drugiCinilac.d4 == true ? true : false}
							/>

							<input
							onChange={(e) => {
								if (drugiCinilac.d5 == true) {
									setDrugiCinilac(drugiCinilac => ({
										...drugiCinilac,
										...{d5: false}
									}))
								} else {
									setDrugiCinilac(drugiCinilac => ({
										...drugiCinilac,
										...{d5: true}
									}))
								}
							}}
							className="w-fit"
							type="checkbox"
							id="d5"
							name="d"
							value="5"
							checked={drugiCinilac.d5 == true ? true : false}
							/>

							<input
							onChange={(e) => {
								if (drugiCinilac.d6 == true) {
									setDrugiCinilac(drugiCinilac => ({
										...drugiCinilac,
										...{d6: false}
									}))
								} else {
									setDrugiCinilac(drugiCinilac => ({
										...drugiCinilac,
										...{d6: true}
									}))
								}
							}}
							className="w-fit"
							type="checkbox"
							id="d6"
							name="d"
							value="6"
							checked={drugiCinilac.d6 == true ? true : false}
							/>

							<input
							onChange={(e) => {
								if (drugiCinilac.d7 == true) {
									setDrugiCinilac(drugiCinilac => ({
										...drugiCinilac,
										...{d7: false}
									}))
								} else {
									setDrugiCinilac(drugiCinilac => ({
										...drugiCinilac,
										...{d7: true}
									}))
								}
							}}
							className="w-fit"
							type="checkbox"
							id="d7"
							name="d"
							value="7"
							checked={drugiCinilac.d7 == true ? true : false}
							/>

							<input
							onChange={(e) => {
								if (drugiCinilac.d8 == true) {
									setDrugiCinilac(drugiCinilac => ({
										...drugiCinilac,
										...{d8: false}
									}))
								} else {
									setDrugiCinilac(drugiCinilac => ({
										...drugiCinilac,
										...{d8: true}
									}))
								}
							}}
							className="w-fit"
							type="checkbox"
							id="d8"
							name="d"
							value="8"
							checked={drugiCinilac.d8 == true ? true : false}
							/>

							<input
							onChange={(e) => {
								if (drugiCinilac.d9 == true) {
									setDrugiCinilac(drugiCinilac => ({
										...drugiCinilac,
										...{d9: false}
									}))
								} else {
									setDrugiCinilac(drugiCinilac => ({
										...drugiCinilac,
										...{d9: true}
									}))
								}
							}}
							className="w-fit"
							type="checkbox"
							id="d9"
							name="d"
							value="9"
							checked={drugiCinilac.d9 == true ? true : false}
							/>

						</div>

					</div>


					<button onClick={vidiOdgovor}
					className="ml-10 my-2 relative
					border-2 border-gray-500 rounded-md p-1 bg-lime-50 hover:bg-lime-200" >
						{props.i18n.Result}</button>

					<button onClick={novoPitanje}
					className="ml-8 my-2 relative
					border-2 border-gray-500 rounded-md p-1 bg-blue-300 hover:bg-blue-400">
						{props.i18n.New_a}</button>
					<br/>

					<O2m_mnozenje
						novo_pitanje={novo_pitanje}
						vidljiv_odgovor={vidljiv_odgovor}
						setOdgovor={setOdgovor}
						rezultat={rezultat}
						resenje={resenje}
					/>
				</>
				)}

			</div>
			<br/>
		</>
	)
}


// const container = document.getElementById('root');
// const root = ReactDOM.createRoot(container);
// root.render(<Zadaci_o2 />);
